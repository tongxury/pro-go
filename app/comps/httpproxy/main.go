package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	// 定义默认的超时和限制
	defaultPort        = "6060"
	maxConcurrentConns = 100                // 最大并发连接数
	dialTimeout        = 1000 * time.Second // 连接目标服务器的超时时间
	idleTimeout        = 1200 * time.Second // 连接空闲超时时间
	// 默认认证信息（生产环境应从环境变量读取）
	defaultUsername = "proxy"
	defaultPassword = "strOngPAssWOrd"
)

// acls 包含访问控制配置
type acls struct {
	allowedIPs []*net.IPNet // IP白名单
	username   string       // 认证用户名
	password   string       // 认证密码
}

// SecureProxy 是一个处理 CONNECT 请求的安全处理器
type SecureProxy struct {
	logger    *log.Logger
	acls      *acls
	semaphore chan struct{} // 用于限制并发连接的信号量
}

// NewSecureProxy 创建并初始化一个 SecureProxy 实例
func NewSecureProxy() *SecureProxy {
	// === 安全修复 1: 允许所有IP访问 ===
	allowedRanges := []string{
		"0.0.0.0/0", // 允许所有IPv4地址
		"::/0",      // 允许所有IPv6地址
	}
	var allowedIPs []*net.IPNet
	for _, cidr := range allowedRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatalf("无法解析白名单IP范围 %s: %v", cidr, err)
		}
		allowedIPs = append(allowedIPs, ipNet)
	}

	// 从环境变量读取认证信息
	username := os.Getenv("PROXY_USERNAME")
	if username == "" {
		username = defaultUsername
	}
	password := os.Getenv("PROXY_PASSWORD")
	if password == "" {
		password = defaultPassword
	}

	return &SecureProxy{
		logger: log.New(os.Stdout, "[Proxy] ", log.LstdFlags),
		acls: &acls{
			allowedIPs: allowedIPs,
			username:   username,
			password:   password,
		},
		// === 安全修复 2: 并发连接数限制 ===
		semaphore: make(chan struct{}, maxConcurrentConns),
	}
}

// authenticateRequest 验证HTTP基本认证
func (p *SecureProxy) authenticateRequest(r *http.Request) bool {
	// 获取Authorization头
	auth := r.Header.Get("Authorization")
	if auth == "" {
		// 尝试从Proxy-Authorization头获取（标准的代理认证方式）
		auth = r.Header.Get("Proxy-Authorization")
	}

	if auth == "" {
		p.logger.Printf("缺少认证头")
		return false
	}

	// 检查是否是Basic认证
	if !strings.HasPrefix(auth, "Basic ") {
		p.logger.Printf("不支持的认证类型")
		return false
	}

	// 解码Base64编码的用户名:密码
	encoded := strings.TrimPrefix(auth, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		p.logger.Printf("无法解码认证信息: %v", err)
		return false
	}

	// 分割用户名和密码
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		p.logger.Printf("认证格式错误")
		return false
	}

	username, password := credentials[0], credentials[1]
	if username == p.acls.username && password == p.acls.password {
		return true
	}

	p.logger.Printf("认证失败: 用户名=%s", username)
	return false
}

// ServeHTTP 实现 http.Handler 接口，处理代理请求
func (p *SecureProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 尝试获取一个信号量，如果满了直接拒绝服务
	select {
	case p.semaphore <- struct{}{}:
		defer func() { <-p.semaphore }() // 请求结束时释放信号量
	default:
		p.logger.Printf("并发连接数已达上限，拒绝请求")
		http.Error(w, "Too many concurrent connections", http.StatusServiceUnavailable)
		return
	}

	// === 新增：密码认证 ===
	if !p.authenticateRequest(r) {
		// 返回407 Proxy Authentication Required
		w.Header().Set("Proxy-Authenticate", "Basic realm=\"Secure Proxy\"")
		http.Error(w, "Proxy Authentication Required", http.StatusProxyAuthRequired)
		return
	}

	// 检查客户端IP是否在白名单内
	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		p.logger.Printf("无法解析客户端地址: %v", err)
		http.Error(w, "Cannot parse client address", http.StatusBadRequest)
		return
	}
	if !p.isIPAllowed(net.ParseIP(clientIP)) {
		p.logger.Printf("拒绝来自未经授权IP的连接: %s", clientIP)
		http.Error(w, "IP not allowed", http.StatusForbidden)
		return
	}

	// 代理只处理 CONNECT 请求
	if r.Method != http.MethodConnect {
		p.logger.Printf("不支持的请求方法: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// === 安全修复 3: 防止SSRF (TOCTOU) 攻击 ===
	// 1. 解析并验证目标地址，获取一个可用的公网IP
	resolvedIP, err := p.validateTarget(r.Host)
	if err != nil {
		p.logger.Printf("非法的目标地址 %s: %v", r.Host, err)
		http.Error(w, "Invalid target host", http.StatusBadRequest)
		return
	}

	// 2. 从原始请求中获取端口号
	_, port, err := net.SplitHostPort(r.Host)
	if err != nil {
		// CONNECT 请求的目标必须是 "host:port" 格式
		p.logger.Printf("无效的目标格式 (必须是 host:port): %s", r.Host)
		http.Error(w, "Invalid target format", http.StatusBadRequest)
		return
	}

	// 3. 使用验证过的IP和原始端口建立连接，避免二次DNS解析
	targetAddr := net.JoinHostPort(resolvedIP, port)
	p.logger.Printf("接受来自 %s 的认证CONNECT请求，目标: %s (解析为 %s)", clientIP, r.Host, targetAddr)

	// 与目标服务器建立连接
	destConn, err := net.DialTimeout("tcp", targetAddr, dialTimeout)
	if err != nil {
		p.logger.Printf("连接目标 %s 失败: %v", targetAddr, err)
		http.Error(w, "Failed to connect to target", http.StatusBadGateway)
		return
	}

	// === 安全修复 4: 正确的 Hijacking 流程 ===
	// 1. 先劫持连接
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		p.logger.Println("服务器不支持 Hijacking")
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		destConn.Close()
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		p.logger.Printf("劫持连接失败: %v", err)
		// 此时可能已经无法发送HTTP错误，但还是尝试一下
		http.Error(w, "Failed to hijack connection", http.StatusInternalServerError)
		destConn.Close()
		return
	}

	// 2. 劫持成功后，手动发送成功响应
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		p.logger.Printf("向客户端发送成功响应失败: %v", err)
		clientConn.Close()
		destConn.Close()
		return
	}

	// === 安全修复 5: Goroutine 和连接生命周期管理 ===
	var once sync.Once
	closeConnections := func() {
		clientConn.Close()
		destConn.Close()
	}

	go p.transfer(destConn, clientConn, &once, closeConnections)
	go p.transfer(clientConn, destConn, &once, closeConnections)
}

// transfer 将数据从源复制到目标，并管理连接超时和关闭
func (p *SecureProxy) transfer(destination net.Conn, source net.Conn, once *sync.Once, closeFunc func()) {
	defer once.Do(closeFunc)

	buffer := make([]byte, 32*1024)
	for {
		source.SetReadDeadline(time.Now().Add(idleTimeout))
		bytesRead, err := source.Read(buffer)
		if bytesRead > 0 {
			destination.SetWriteDeadline(time.Now().Add(idleTimeout))
			_, writeErr := destination.Write(buffer[:bytesRead])
			if writeErr != nil {
				// === 健壮性修复: 使用 errors.Is 进行错误检查 ===
				if !errors.Is(writeErr, net.ErrClosed) {
					p.logger.Printf("写入数据失败: %v", writeErr)
				}
				return
			}
		}
		if err != nil {
			// === 健壮性修复: 使用 errors.Is 进行错误检查 ===
			if err != io.EOF && !errors.Is(err, net.ErrClosed) {
				p.logger.Printf("读取数据失败: %v", err)
			}
			return
		}
	}
}

// isIPAllowed 检查给定的IP是否在白名单内
func (p *SecureProxy) isIPAllowed(ip net.IP) bool {
	if ip == nil {
		return false
	}
	for _, ipNet := range p.acls.allowedIPs {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

// validateTarget 验证目标主机是否为私有或回环地址, 并返回一个可用的公网IP
func (p *SecureProxy) validateTarget(host string) (string, error) {
	hostWithoutPort, _, err := net.SplitHostPort(host)
	if err != nil {
		// 假定没有端口，整个字符串都是主机名
		hostWithoutPort = host
	}

	ips, err := net.LookupIP(hostWithoutPort)
	if err != nil {
		return "", fmt.Errorf("无法解析主机: %v", err)
	}

	var firstPublicIP net.IP
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() || ip.IsMulticast() {
			return "", fmt.Errorf("目标地址解析为私有、回环或特殊用途地址: %s", ip)
		}
		if ip.IsGlobalUnicast() && firstPublicIP == nil {
			firstPublicIP = ip
		}
	}

	if firstPublicIP != nil {
		return firstPublicIP.String(), nil
	}

	return "", fmt.Errorf("未找到任何可用的公网IP地址")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	proxy := NewSecureProxy()

	// 输出认证信息（生产环境中应该避免在日志中显示）
	proxy.logger.Printf("代理认证: 用户名=%s, 密码=%s", proxy.acls.username, "***")

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      proxy,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// === 安全修复 6: 优雅关闭 ===
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		proxy.logger.Printf("启动安全 CONNECT 代理于端口 %s", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			proxy.logger.Fatalf("启动服务器失败: %v", err)
		}
	}()

	<-stopChan // 等待中断信号

	proxy.logger.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		proxy.logger.Fatalf("服务器优雅关闭失败: %v", err)
	}

	proxy.logger.Println("服务器已优雅关闭")
}

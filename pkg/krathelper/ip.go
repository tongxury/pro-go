package krathelper

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport/http"
	"net"
	"strings"
)

func ClientCountryCode(ctx context.Context) string {
	r, ok := http.RequestFromServerContext(ctx)
	if !ok {
		return ""
	}

	return r.Header.Get("Cf-Ipcountry")
}

func ClientPublicIP(ctx context.Context) string {

	r, ok := http.RequestFromServerContext(ctx)
	if !ok {
		return ""
	}

	//log.Debugw("ClientPublicIP", r.Header)

	//Cf-Ipcountry cloudflare 可以直接获取地区
	ips := r.Header.Get("Cf-Connecting-Ip")
	if ip := findPublicIP(ips); ip != "" {
		return ip
	}

	ips = r.Header.Get("X-Original-Forwarded-For")
	if ip := findPublicIP(ips); ip != "" {
		return ip
	}

	ips = r.Header.Get("X-Forwarded-For")
	if ip := findPublicIP(ips); ip != "" {
		return ip
	}

	ips = r.Header.Get("X-Real-Ip")
	if ip := findPublicIP(ips); ip != "" {
		return ip
	}

	ips = RemoteIP(r)
	if ip := findPublicIP(ips); ip != "" {
		return ip
	}

	return ""
}

func findPublicIP(ips string) string {

	if ips == "" {
		return ""
	}

	for _, x := range strings.Split(ips, ",") {
		if !isLocalIP(x) {
			return strings.TrimSpace(x)
		}
	}
	return ""
}

// HasLocalIPAddr 检测 IP 地址字符串是否是内网地址
func isLocalIP(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP 检测 IP 地址是否是内网地址
// 通过直接对比ip段范围效率更高
func HasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}

// RemoteIP 通过 RemoteAddr 获取 IP 地址， 只是一个快速解析方法。
func RemoteIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

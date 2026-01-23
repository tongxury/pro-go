package google

import (
	"context"
	"encoding/base64"
	"strings"

	"google.golang.org/api/gmail/v1"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SendEmail(ctx context.Context) error {

	return nil
}

func createMessage(from, to, subject, bodyContent string) *gmail.Message {
	// 构建邮件头部和内容
	emailLines := []string{
		"From: " + from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		bodyContent,
	}

	// 合并成完整邮件
	emailContent := strings.Join(emailLines, "\r\n")

	// Base64 URL编码
	encodedMessage := base64.URLEncoding.EncodeToString([]byte(emailContent))

	// 创建Gmail消息对象
	return &gmail.Message{Raw: encodedMessage}
}

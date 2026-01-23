package mail

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

type Client struct {
	conf Config
}

type Config struct {
	Host     string
	Port     string
	From     string
	Password string
}

func NewClient(conf Config) *Client {

	//conf.Host = "108.181.252.240"
	//conf.Port = "587"
	//conf.From = "noreply@mail.veogo.ai"
	//conf.Password = "1qaz@WSXveogo"
	//587
	conf.Host = "smtpdm.aliyun.com"
	conf.Port = "80"
	conf.From = "noreply@mail.veogo.ai"
	conf.Password = "1qaz2WSXveogo"

	//1qaz2WSXveogo

	return &Client{
		conf: conf,
	}
}

type Params struct {
	Subject string
	To      string
	Content string
}

func (t *Client) Send(params Params) error {

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", params.Subject, t.conf.From)
	e.To = []string{params.To}
	e.Subject = params.Subject
	e.HTML = []byte(params.Content)

	//host := "138.201.193.45"

	err := e.SendWithStartTLS(t.conf.Host+":"+t.conf.Port, smtp.PlainAuth("",
		t.conf.From, t.conf.Password, t.conf.Host),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return err
	}

	return nil
}

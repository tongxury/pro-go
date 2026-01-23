package awsmail

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Client struct {
	*ses.SES
}

func NewClient(conf Config) *Client {

	region := "us-east-1"

	if conf.Region != "" {
		region = conf.Region
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(conf.AccessKey, conf.AccessSecret, ""),
		Region:      aws.String(region),
		//DisableSSL:  aws.Bool(true),
	})

	if err != nil {
		panic(err)
	}
	return &Client{
		SES: ses.New(sess),
	}
}

func (t *Client) Send(sender, recipient, subject, htmlTpl string) error {

	var charSet = "UTF-8"

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(htmlTpl),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String("TextBody"),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	_, err := t.SendEmail(input)

	return err
}

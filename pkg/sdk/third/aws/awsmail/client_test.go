package awsmail

import (
	"store/pkg/sdk/third/aws"
	"strings"
	"testing"

	"store/confs"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "noreply@veogo.ai"

	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	Recipient = "tongxurt@gmail.com"
	//Recipient = "sifan536414@gmail.com"
	//Recipient = "2491764611marcus@gmail.com"
	//Recipient = "536414887@qq.com"
	//Recipient = "2491764611@qq.com"

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The subject line for the email.
	Subject = "Amazon SES Test (AWS SDK for Go)"
	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	htmlBody = `
<!DOCTYPE html>
		<html lang="en" xmlns:th="http://www.thymeleaf.org">
		<head>
			<meta charset="UTF-8">
			<title></title>
			<style>
				table {
					width: 700px;
					margin: 0 auto;
				}
		
				#top {
					width: 700px;
					border-bottom: 1px solid #ccc;
					margin: 0 auto 30px;
				}
		
				#top table {
					font: 12px Tahoma, Arial, 宋体;
					height: 40px;
				}
		
				#content {
					width: 680px;
					padding: 0 10px;
					margin: 0 auto;
				}
		
				#content_top {
					line-height: 1.5;
					font-size: 14px;
					margin-bottom: 25px;
					color: #4d4d4d;
				}
		
				#content_top strong {
					display: block;
					margin-bottom: 15px;
				}
		
				#content_top strong span {
					color: #f60;
					font-size: 16px;
				}
		
				#verificationCode {
					color: #f60;
					font-size: 24px;
				}
		
				#content_bottom {
					margin-bottom: 30px;
				}
		
				#content_bottom small {
					display: block;
					margin-bottom: 20px;
					font-size: 12px;
					color: #747474;
				}
		
				#bottom {
					width: 700px;
					margin: 0 auto;
				}
		
				#bottom div {
					padding: 10px 10px 0;
					border-top: 1px solid #ccc;
					color: #747474;
					margin-bottom: 20px;
					line-height: 1.3em;
					font-size: 12px;
				}
		
				#content_top strong span {
					font-size: 18px;
					color: #FE4F70;
				}
		
				#sign {
					text-align: right;
					font-size: 18px;
					color: #FE4F70;
					font-weight: bold;
				}
		
				#verificationCode {
					height: 100px;
					width: 680px;
					text-align: center;
					margin: 30px 0;
				}
		
				#verificationCode div {
					height: 100px;
					width: 680px;
		
				}
		
				.button {
					color: #FE4F70;
					margin-left: 10px;
					height: 80px;
					width: 80px;
					resize: none;
					font-size: 42px;
					border: none;
					outline: none;
					padding: 10px 15px;
					background: #ededed;
					text-align: center;
					border-radius: 17px;
					box-shadow: 6px 6px 12px #cccccc,
					-6px -6px 12px #ffffff;
				}
		
				.button:hover {
					box-shadow: inset 6px 6px 4px #d1d1d1,
					inset -6px -6px 4px #ffffff;
				}
		
			</style>
		</head>
		<body>
		<table>
			<tbody>
			<tr>
				<td>
					<div id="top">
						<table>
							<tbody>
							<tr>
								<td></td>
							</tr>
							</tbody>
						</table>
					</div>
		
					<div id="content">
						<div id="content_top">
		
							<strong>Dear <span>__USER_NAME__</span></strong>
							<strong>
								Thank you for signing up for our service. To complete the registration process and activate your
								account, please use the following verification code:
							</strong>
							<div id="verificationCode">
								Verification Code:<span>__VERIFY_CODE__</span>
							</div>
						</div>
						<div id="content_bottom">
							<small>
								Please enter this code on the registration page to verify your account. If you did not sign up
								for our service, please disregard this email.
								<br>
								If you have any questions or need assistance, please feel free to contact our support team at
								service@xbuddy.ai.
							</small>
						</div>
					</div>
					<div id="bottom">
						<div>
							<p>Thank you for choosing our service.<br>
								Best regards,
								StudyGPT team
							</p>
						</div>
					</div>
				</td>
			</tr>
			</tbody>
		</table>
		</body>
`
	// The character encoding for the email.
	CharSet = "UTF-8"
)

func TestClient_Send(t *testing.T) {

	body := strings.ReplaceAll(htmlBody, "__USER_NAME__", "赵斯凡")
	body = strings.ReplaceAll(body, "__GUIDE_URL__", "https://bc9fx7uhk6.feishu.cn/docx/QAowdSUZkoLvEAx71KbcuMqznFd")
	body = strings.ReplaceAll(body, "__DISCORD_URL__", "https://discord.com/invite/cqgndY8vGA")

	err := NewClient(Config{
		AwsConfig: aws.AwsConfig{
			AccessKey:    confs.AWSAccessKey,
			AccessSecret: confs.AWSAccessSecret,
			//AccessKey:    confs.AWSAccessKey,
			//AccessSecret: confs.AWSAccessSecret,
			Region: "us-east-1",
		},
	}).Send(Sender, Recipient, Subject, body)
	if err != nil {
		return
	}
}

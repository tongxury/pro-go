package awsmail

import "store/pkg/sdk/third/aws"

type Config struct {
	aws.AwsConfig
	Sender string
}

package form

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/aws"
)

var s3uploader *s3manager.Uploader

func SetupS3() {
	s := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	}))
	s3uploader = s3manager.NewUploader(s)
}

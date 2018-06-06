package form

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var s3uploader *s3manager.Uploader

func SetupS3() {
	s := session.Must(session.NewSession())
	s3uploader = s3manager.NewUploader(s)
}

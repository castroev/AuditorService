package s3service

import (
	"bytes"
	"time"

	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/config"
	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/customclient"
	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var s3session *session.Session

//SaveToS3Bucket saves the audit log json in an s3 bucket
func SaveToS3Bucket(key string, data []byte) bool {
	c := config.GetConfig()
	logging.Logger.Infof("Saving %s to %s.", key, c.S3.Bucket)

	svc := s3manager.NewUploader(s3session)

	_, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.S3.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		logging.Logger.Errorf("Error saving to S3 key (%s) %v", key, err)
		return false
	}
	return true
}

// CreateS3Session creates a global s3 session
func CreateS3Session() error {
	c := config.GetConfig()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.S3.Region),
		HTTPClient: customclient.NewHTTPClientWithTimeouts(customclient.ClientSettings{
			Connect:          5 * time.Second,
			ExpectContinue:   1 * time.Second,
			IdleConn:         90 * time.Second,
			ConnKeepAlive:    90 * time.Second,
			MaxAllIdleConns:  100,
			MaxHostIdleConns: 100,
			ResponseHeader:   5 * time.Second,
			TLSHandshake:     5 * time.Second,
		}),
	})
	if err != nil {
		logging.Logger.Errorf("Error creating session for S3: %v", err)
		return nil
	}
	s3session = sess
	return nil
}

func getS3Session() *session.Session {
	return s3session
}

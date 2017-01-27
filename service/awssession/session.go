package awssession

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

// A AWSSession provides a aws client session
type AWSSession struct {
	Session *session.Session
	Config  *aws.Config
}

// New returns a AWSSession
func New() (*AWSSession, error) {

	region := ""
	region = "ap-northeast-1"

	config := aws.NewConfig()
	config.WithRegion(region)

	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	return &AWSSession{
		Session: sess,
		Config:  config,
	}, nil
}

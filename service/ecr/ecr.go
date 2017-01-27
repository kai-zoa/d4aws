package ecr

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/kai-zoa/d4aws/service/awssession"
	"github.com/pkg/errors"
	"strings"
)

// An ECR provides an ECR client session
type ECR struct {
	*awssession.AWSSession
}

// A GetLoginCommandInput is GetLoginCommand's input parameters.
type GetLoginCommandInput struct {
	RegistryIDs []string
}

// New returns a ECR client session
func New() (*ECR, error) {
	session, err := awssession.New()
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	return &ECR{
		AWSSession: session,
	}, nil
}

// GetLoginCommand returns a Docker LoginCommand for login to ECR
func (s *ECR) GetLoginCommand(input *GetLoginCommandInput) (string, error) {
	sess, config := s.Session, s.Config

	var ids []*string

	if len(input.RegistryIDs) > 0 {
		ids = make([]*string, len(input.RegistryIDs))
		for i := 0; i < len(ids); i++ {
			ids[i] = &input.RegistryIDs[i]
		}
	}
	output, err := ecr.New(sess, config).GetAuthorizationToken(
		&ecr.GetAuthorizationTokenInput{RegistryIds: ids},
	)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	if len(output.AuthorizationData) == 0 {
		return "", errors.New("failed0")
	}

	registry := *(output.AuthorizationData[0].ProxyEndpoint)

	tokenData, err := base64.StdEncoding.DecodeString(*(output.AuthorizationData[0].AuthorizationToken))
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	pair := strings.Split(string(tokenData), ":")

	if len(pair) != 2 {
		return "", errors.New("failed1")
	}

	return fmt.Sprintf("docker login -u %s -p %s %s", pair[0], pair[1], registry), nil
}

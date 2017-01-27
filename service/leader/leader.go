package leader

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kai-zoa/d4aws/service/awssession"
	"github.com/pkg/errors"
)

type Leader struct {
	*awssession.AWSSession
	StackName string
}

func New(stackName string) (*Leader, error) {
	session, err := awssession.New()
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	return &Leader{
		AWSSession: session,
		StackName:  stackName,
	}, nil
}

func (s *Leader) GetPrivateIPAddress() (string, error) {
	sess, config := s.Session, s.Config

	tableName, err := s.getPhysicalResouceID("SwarmDynDBTable")
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	scanOutput, err := dynamodb.New(sess, config).Scan(
		&dynamodb.ScanInput{
			TableName: &tableName,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	if len(scanOutput.Items) == 0 {
		return "", errors.New("failed01")
	}

	privateIP := ""
	for _, e := range scanOutput.Items {
		if v, ok := e["node_type"]; !ok || (*v.S) != "primary_manager" {
			continue
		}
		if v, ok := e["ip"]; ok {
			privateIP = (*v.S)
			break
		}
	}

	if privateIP == "" {
		return "", errors.New("failed02")
	}

	return privateIP, nil
}

func (s *Leader) GetPublicIPAddress() (string, error) {

	privateIP, err := s.GetPrivateIPAddress()
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	sess, config := s.Session, s.Config

	managerAsgID, err := s.getPhysicalResouceID("ManagerAsg")
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	groupsOutput, err := autoscaling.New(sess, config).DescribeAutoScalingGroups(
		&autoscaling.DescribeAutoScalingGroupsInput{
			AutoScalingGroupNames: []*string{&managerAsgID},
		},
	)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	if len(groupsOutput.AutoScalingGroups) == 0 {
		return "", errors.New("failed02")
	}

	ids := []*string{}
	for _, ins := range groupsOutput.AutoScalingGroups[0].Instances {
		ids = append(ids, ins.InstanceId)
	}

	insOutput, err := ec2.New(sess, config).DescribeInstances(
		&ec2.DescribeInstancesInput{
			InstanceIds: ids,
		},
	)

	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}

	for _, res := range insOutput.Reservations {
		for _, ins := range res.Instances {
			if *ins.PrivateIpAddress == privateIP {
				return *ins.PublicIpAddress, nil
			}
		}
	}

	return "", errors.New("instance not found")
}

func (s *Leader) getPhysicalResouceID(resKey string) (string, error) {
	sess, config, stackName := s.Session, s.Config, s.StackName

	r, err := cloudformation.New(sess, config).DescribeStackResource(
		&cloudformation.DescribeStackResourceInput{
			StackName:         &stackName,
			LogicalResourceId: &resKey,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, err.Error())
	}
	return *r.StackResourceDetail.PhysicalResourceId, nil
}

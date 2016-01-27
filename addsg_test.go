package main

import (
	"fmt"
	"testing"

	"github.com/amanya/addsg/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AddSGSuite struct {
	suite.Suite
}

func TestAddSGSuite(t *testing.T) {
	suite.Run(t, new(AddSGSuite))
}

func (s *AddSGSuite) TestGetSecurityGroupNotExists() {
	client := new(mocks.EC2er)
	client.On("DescribeSecurityGroups", mock.AnythingOfType("*ec2.DescribeSecurityGroupsInput")).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

	helper := &EC2Helper{client}

	sg, err := helper.getSecurityGroup("asdf", "zxcv")
	s.Equal(sg, "")
	s.Equal(err, nil)
}

func (s *AddSGSuite) TestGetSecurityGroupDidExists() {
	client := new(mocks.EC2er)
	client.On("DescribeSecurityGroups", mock.AnythingOfType("*ec2.DescribeSecurityGroupsInput")).Return(
		&ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: []*ec2.SecurityGroup{
				&ec2.SecurityGroup{
					GroupId: aws.String("asdf"),
				},
			},
		}, nil)

	helper := &EC2Helper{client}

	sg, err := helper.getSecurityGroup("asdf", "zxcv")
	s.Equal(sg, "asdf")
	s.Equal(err, nil)
}

func (s *AddSGSuite) TestGetInstanceNotFound() {
	client := new(mocks.EC2er)
	client.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(
		&ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{},
		}, nil)

	helper := &EC2Helper{client}

	i, err := helper.getInstance("1.1.1.1")
	assert.Nil(s.T(), i)
	s.Equal(err, fmt.Errorf("instance not found"))
}

func (s *AddSGSuite) TestGetInstanceFound() {
	client := new(mocks.EC2er)
	client.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(
		&ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{
				&ec2.Reservation{
					Instances: []*ec2.Instance{
						&ec2.Instance{
							InstanceId: aws.String("i-asdfasdf"),
						},
					},
				},
			},
		}, nil)

	helper := &EC2Helper{client}

	i, err := helper.getInstance("1.1.1.1")
	s.Equal(i.InstanceId, aws.String("i-asdfasdf"))
	s.Equal(err, nil)
}

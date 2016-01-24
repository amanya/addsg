package main

import (
	"github.com/amanya/addsg/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
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

package main

import (
	"errors"
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

func (s *AddSGSuite) TestCreateSecurityGroup() {
	client := new(mocks.EC2er)
	client.On("CreateSecurityGroup", mock.AnythingOfType("*ec2.CreateSecurityGroupInput")).Return(
		&ec2.CreateSecurityGroupOutput{
			GroupId: aws.String("sg-asdf"),
		}, nil)

	helper := &EC2Helper{client}

	sg, err := helper.createSecurityGroup("vpc-asdf", "asdf")

	s.Equal(sg, "sg-asdf")
	s.Equal(err, nil)
}

func (s *AddSGSuite) TestAddIpToSecurityGroup() {
	client := new(mocks.EC2er)
	client.On("AuthorizeSecurityGroupIngress", mock.AnythingOfType("*ec2.AuthorizeSecurityGroupIngressInput")).Return(nil, nil)

	helper := &EC2Helper{client}

	err := helper.addIpToSecurityGroup("1.2.3.4/32", "sg-asdf")
	s.Equal(err, nil)
}

func (s *AddSGSuite) TestInstanceAlreadyInSecurityGroup() {
	client := new(mocks.EC2er)

	i := &ec2.Instance{
		SecurityGroups: []*ec2.GroupIdentifier{
			&ec2.GroupIdentifier{
				GroupId: aws.String("sg-asdf"),
			},
		},
	}

	helper := &EC2Helper{client}

	err := helper.addSecurityGroupToInstance(i, "sg-asdf")
	s.Equal(err, errors.New("instance already has the security group"))
}

func (s *AddSGSuite) TestAddSecurityGroupToInstance() {
	client := new(mocks.EC2er)

	i := &ec2.Instance{
		InstanceId: aws.String("i-asdf"),
		SecurityGroups: []*ec2.GroupIdentifier{
			&ec2.GroupIdentifier{
				GroupId: aws.String("sg-asdf"),
			},
		},
	}

	client.On("ModifyInstanceAttribute", mock.AnythingOfType("*ec2.ModifyInstanceAttributeInput")).Return(nil, nil)

	helper := &EC2Helper{client}

	err := helper.addSecurityGroupToInstance(i, "sg-qwer")
	s.Equal(err, nil)
}

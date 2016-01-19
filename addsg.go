package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/rdegges/go-ipify"
	"log"
	"os"
)

type EC2er interface {
	DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	CreateSecurityGroup(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error)
	AuthorizeSecurityGroupIngress(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
	ModifyInstanceAttribute(*ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error)
}

var _ EC2er = (*ec2.EC2)(nil)

var (
	instanceIp = flag.String("i", "", "IP of the instance we want to access")
)

func usage() {
	log.Printf("usage: addsg -i [instance-ip]")
	os.Exit(1)
}

func getSecurityGroup(e EC2er, vpcId string, sgName string) (string, error) {
	r, err := e.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcId),
				},
			},
			&ec2.Filter{
				Name: aws.String("group-name"),
				Values: []*string{
					aws.String(sgName),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	if len(r.SecurityGroups) == 0 {
		return "", nil
	}

	return *r.SecurityGroups[0].GroupId, nil
}

func getInstance(e EC2er, instanceIp string) (*ec2.Instance, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("ip-address"),
				Values: []*string{
					aws.String(instanceIp),
				},
			},
		},
	}

	r, err := e.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	if len(r.Reservations) == 0 {
		return nil, fmt.Errorf("instance not found")
	}

	return r.Reservations[0].Instances[0], nil
}

func createSecurityGroup(e EC2er, vpcId string, sgName string) (string, error) {
	securityGroupOpts := &ec2.CreateSecurityGroupInput{}
	securityGroupOpts.VpcId = aws.String(vpcId)
	securityGroupOpts.Description = aws.String("Created by addsg")
	securityGroupOpts.GroupName = aws.String(sgName)

	r, err := e.CreateSecurityGroup(securityGroupOpts)
	if err != nil {
		return "", err
	}

	return *r.GroupId, nil
}

func addIpToSecurityGroup(e EC2er, ip string, sgId string) error {
	r, err := e.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:     aws.String(ip + "/32"),
		FromPort:   aws.Int64(22),
		ToPort:     aws.Int64(22),
		GroupId:    aws.String(sgId),
		IpProtocol: aws.String("TCP"),
	})
	_ = r

	return err
}

func addSecurityGroupToInstance(e EC2er, i *ec2.Instance, sgId string) (string, error) {
	var groups []*string
	for _, group := range i.SecurityGroups {
		if *group.GroupId != sgId {
			groups = append(groups, group.GroupId)
		} else {
			return "", errors.New("instance already has the security group")
		}
	}

	groups = append(groups, &sgId)

	_, err := e.ModifyInstanceAttribute(&ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(*i.InstanceId),
		Groups:     groups,
	})
	if err != nil {
		return "", err
	}
	return "", nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *instanceIp == "" {
		usage()
	}

	sess := session.New()
	var e EC2er = ec2.New(sess)

	i, err := getInstance(e.(*ec2.EC2), *instanceIp)
	if err != nil {
		log.Printf("Could't find the instance: %s", err)
		os.Exit(1)
	}

	log.Printf("Found instance id: %+v", i.InstanceId)
	log.Printf("Found VPC id: %s", *i.VpcId)

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting the hostname: %s", err)
		os.Exit(1)
	}
	sgName := "addsg-" + hostname

	sgId, err := getSecurityGroup(e.(*ec2.EC2), *i.VpcId, sgName)
	if err != nil {
		log.Printf("Error searching the sg: %s", err)
		os.Exit(1)
	}

	if sgId == "" {
		// create the sg
		sg, err := createSecurityGroup(e.(*ec2.EC2), *i.VpcId, sgName)
		sgId = sg
		if err != nil {
			log.Printf("Error creating the sg: %s", err)
			os.Exit(1)
		}

		log.Printf("Created sg: %s", sgId)
	} else {
		log.Printf("Recycling sg: %s", sgId)
	}

	ip, err := ipify.GetIp()
	if err != nil {
		log.Printf("Couldn't get IP address: %s", err)
		os.Exit(1)
	}
	log.Printf("Granting access to IP: %s", ip)

	err = addIpToSecurityGroup(e.(*ec2.EC2), ip, sgId)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "InvalidPermission.Duplicate" {
				log.Printf("The IP was already on the sg")
			} else {
				log.Printf("Couldn't add IP to sg: %s", err)
				os.Exit(1)
			}
		} else {
			log.Printf("Couldn't add IP to sg: %s", err)
			os.Exit(1)
		}
	}

	s, err := addSecurityGroupToInstance(e.(*ec2.EC2), i, sgId)
	_ = s
	if err != nil {
		log.Printf("Couldn't add instance to sg: %s", err)
		os.Exit(1)
	}
	log.Printf("Done")

}

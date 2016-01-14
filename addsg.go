package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"os"
)

var (
	instanceIp = flag.String("i", "", "IP of the instance we want to access")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: addsg -i [instance-ip]\n")
	os.Exit(1)
}

func describeSecurityGroups(e *ec2.EC2) (string, error) {
	params := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("group-name"),
				Values: []*string{
					aws.String("addsg"),
				},
			},
		},
	}
	r, err := e.DescribeSecurityGroups(params)
	if err != nil {
		return "", err
	}

	if len(r.SecurityGroups) == 0 {
		return "", fmt.Errorf("security group not found")
	}

	return *r.SecurityGroups[0].GroupId, nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *instanceIp == "" {
		usage()
	}

	sess := session.New()
	e := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("ip-address"),
				Values: []*string{
					aws.String(*instanceIp),
				},
			},
		},
	}

	resp, _ := e.DescribeInstances(params)

	var instanceId = *resp.Reservations[0].Instances[0].InstanceId

	fmt.Fprintf(os.Stdout, "instance: %s\n", instanceId)
	r, err := describeSecurityGroups(e)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err)
	}
	fmt.Fprintf(os.Stdout, "sg: %s\n", r)
}

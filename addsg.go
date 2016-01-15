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

func getSecurityGroup(e *ec2.EC2, sgName string) (string, error) {
	r, err := e.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
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
		return "", fmt.Errorf("security group not found")
	}

	return *r.SecurityGroups[0].GroupId, nil
}

func getInstanceId(e *ec2.EC2, instanceIp string) (string, error) {
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
		return "", err
	}

	if len(r.Reservations) == 0 {
		return "", fmt.Errorf("instance not found")
	}

	return *r.Reservations[0].Instances[0].InstanceId, nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *instanceIp == "" {
		usage()
	}

	sess := session.New()
	e := ec2.New(sess)

	iId, err := getInstanceId(e, *instanceIp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "iId: %s\n", iId)

	sgId, err := getSecurityGroup(e, "addsg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "sgId: %s\n", sgId)
}

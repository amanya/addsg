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

func describeSecurityGroups(svc *ec2.EC2) string {
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
	resp, _ := svc.DescribeSecurityGroups(params)
	if len(resp.SecurityGroups) > 0 {
		return *resp.SecurityGroups[0].GroupId
	}
	return ""
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *instanceIp == "" {
		usage()
	}

	sess := session.New()
	svc := ec2.New(sess)

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

	resp, _ := svc.DescribeInstances(params)

	var instanceId = *resp.Reservations[0].Instances[0].InstanceId

	fmt.Fprintf(os.Stdout, "instance: %s\n", instanceId)
	fmt.Fprintf(os.Stdout, "sg: %s\n", describeSecurityGroups(svc))
}

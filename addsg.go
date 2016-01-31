package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rdegges/go-ipify"
)

type EC2er interface {
	DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	CreateSecurityGroup(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error)
	AuthorizeSecurityGroupIngress(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
	ModifyInstanceAttribute(*ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error)
	DeleteSecurityGroup(*ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error)
}

var _ EC2er = (*ec2.EC2)(nil)

type EC2Helper struct {
	client EC2er
}

var (
	cleanMode  = flag.Bool("c", false, "Remove all security groups created by me")
	instanceIp = flag.String("i", "", "IP of the instance we want to access")
)

func usage() {
	log.Printf("usage: addsg -i [instance-ip]")
	os.Exit(1)
}

func (e *EC2Helper) getSecurityGroup(vpcId string, sgName string) (string, error) {
	r, err := e.client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
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

func (e *EC2Helper) getAllSecurityGroups(sgName string) ([]*ec2.SecurityGroup, error) {
	r, err := e.client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
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
		return nil, err
	}

	return r.SecurityGroups, nil
}

func (e *EC2Helper) deleteSecurityGroup(sg *ec2.SecurityGroup) error {
	_, err := e.client.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
		GroupId: sg.GroupId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2Helper) getInstance(instanceIp string) (*ec2.Instance, error) {
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

	r, err := e.client.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	if len(r.Reservations) == 0 {
		return nil, fmt.Errorf("instance not found")
	}

	return r.Reservations[0].Instances[0], nil
}

func (e *EC2Helper) getAllReservationsByVpcAndSG(group *ec2.SecurityGroup) ([]*ec2.Reservation, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("vpc-id"),
				Values: []*string{
					group.VpcId,
				},
			},
		},
	}

	r, err := e.client.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	return r.Reservations, nil
}

func (e *EC2Helper) createSecurityGroup(vpcId string, sgName string) (string, error) {
	securityGroupOpts := &ec2.CreateSecurityGroupInput{}
	securityGroupOpts.VpcId = aws.String(vpcId)
	securityGroupOpts.Description = aws.String("Created by addsg")
	securityGroupOpts.GroupName = aws.String(sgName)

	r, err := e.client.CreateSecurityGroup(securityGroupOpts)
	if err != nil {
		return "", err
	}

	return *r.GroupId, nil
}

func (e *EC2Helper) addIpToSecurityGroup(ip string, sgId string) error {
	r, err := e.client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:     aws.String(ip + "/32"),
		FromPort:   aws.Int64(22),
		ToPort:     aws.Int64(22),
		GroupId:    aws.String(sgId),
		IpProtocol: aws.String("TCP"),
	})
	_ = r

	return err
}

func (e *EC2Helper) addSecurityGroupToInstance(i *ec2.Instance, sgId string) error {
	var groups []*string
	for _, group := range i.SecurityGroups {
		if *group.GroupId != sgId {
			groups = append(groups, group.GroupId)
		} else {
			return errors.New("instance already has the security group")
		}
	}

	groups = append(groups, &sgId)

	_, err := e.client.ModifyInstanceAttribute(&ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(*i.InstanceId),
		Groups:     groups,
	})

	return err
}

func (e *EC2Helper) removeSecurityGroupFromInstance(i *ec2.Instance, sgId string) (bool, error) {
	var groups []*string
	found := false
	for _, group := range i.SecurityGroups {
		if *group.GroupId != sgId {
			groups = append(groups, group.GroupId)
		} else {
			found = true
		}
	}

	_, err := e.client.ModifyInstanceAttribute(&ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(*i.InstanceId),
		Groups:     groups,
	})
	if err != nil {
		return false, err
	}
	return found, nil
}

func (e *EC2Helper) cleanUp() {
	sgName := makeSgName()
	groups, _ := e.getAllSecurityGroups(sgName)
	for _, group := range groups {
		log.Printf("Removing sg: %s", *group.GroupId)
		res, err := e.getAllReservationsByVpcAndSG(group)
		if err != nil {
			log.Printf("Error getting reservations: %s", err)
			os.Exit(1)
		}
		for _, r := range res {
			for _, i := range r.Instances {
				found, err := e.removeSecurityGroupFromInstance(i, *group.GroupId)
				if err != nil {
					log.Printf("Error removing sg from instance: %s", err)
					os.Exit(1)
				}
				if found {
					log.Printf("Removed sg %s from instance %s", *group.GroupId, *i.InstanceId)
				}
			}
		}
		err1 := e.deleteSecurityGroup(group)
		if err1 != nil {
			log.Printf("Error removing sg: %s", err1)
			os.Exit(1)
		}
		log.Printf("Removed sg %s", *group.GroupId)
	}
}

func makeSgName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting the hostname: %s", err)
		os.Exit(1)
	}
	return "addsg-" + hostname
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if (*instanceIp == "" && !*cleanMode) || (*instanceIp != "" && *cleanMode) {
		usage()
	}

	sess := session.New()
	var e EC2er = ec2.New(sess)

	helper := &EC2Helper{e}

	if *cleanMode {
		helper.cleanUp()
		os.Exit(0)
	}

	i, err := helper.getInstance(*instanceIp)
	if err != nil {
		log.Printf("Could't find the instance: %s", err)
		os.Exit(1)
	}

	log.Printf("Found instance id: %s", *i.InstanceId)
	log.Printf("Found VPC id: %s", *i.VpcId)

	sgName := makeSgName()

	sgId, err := helper.getSecurityGroup(*i.VpcId, sgName)
	if err != nil {
		log.Printf("Error searching the sg: %s", err)
		os.Exit(1)
	}

	if sgId == "" {
		// create the sg
		sg, err := helper.createSecurityGroup(*i.VpcId, sgName)
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

	err = helper.addIpToSecurityGroup(ip, sgId)
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

	err = helper.addSecurityGroupToInstance(i, sgId)
	if err != nil {
		log.Printf("Couldn't add instance to sg: %s", err)
		os.Exit(1)
	}
	log.Printf("Done")
}

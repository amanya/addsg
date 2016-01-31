package mocks

import "github.com/stretchr/testify/mock"

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2er struct {
	mock.Mock
}

// DescribeSecurityGroups provides a mock function with given fields: _a0
func (_m *EC2er) DescribeSecurityGroups(_a0 *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	ret := _m.Called(_a0)

	var r0 *ec2.DescribeSecurityGroupsOutput
	if rf, ok := ret.Get(0).(func(*ec2.DescribeSecurityGroupsInput) *ec2.DescribeSecurityGroupsOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.DescribeSecurityGroupsOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ec2.DescribeSecurityGroupsInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DescribeInstances provides a mock function with given fields: _a0
func (_m *EC2er) DescribeInstances(_a0 *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	ret := _m.Called(_a0)

	var r0 *ec2.DescribeInstancesOutput
	if rf, ok := ret.Get(0).(func(*ec2.DescribeInstancesInput) *ec2.DescribeInstancesOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.DescribeInstancesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ec2.DescribeInstancesInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSecurityGroup provides a mock function with given fields: _a0
func (_m *EC2er) CreateSecurityGroup(_a0 *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
	ret := _m.Called(_a0)

	var r0 *ec2.CreateSecurityGroupOutput
	if rf, ok := ret.Get(0).(func(*ec2.CreateSecurityGroupInput) *ec2.CreateSecurityGroupOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.CreateSecurityGroupOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ec2.CreateSecurityGroupInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeSecurityGroupIngress provides a mock function with given fields: _a0
func (_m *EC2er) AuthorizeSecurityGroupIngress(_a0 *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	ret := _m.Called(_a0)

	var r0 *ec2.AuthorizeSecurityGroupIngressOutput
	if rf, ok := ret.Get(0).(func(*ec2.AuthorizeSecurityGroupIngressInput) *ec2.AuthorizeSecurityGroupIngressOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.AuthorizeSecurityGroupIngressOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ec2.AuthorizeSecurityGroupIngressInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModifyInstanceAttribute provides a mock function with given fields: _a0
func (_m *EC2er) ModifyInstanceAttribute(_a0 *ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error) {
	ret := _m.Called(_a0)

	var r0 *ec2.ModifyInstanceAttributeOutput
	if rf, ok := ret.Get(0).(func(*ec2.ModifyInstanceAttributeInput) *ec2.ModifyInstanceAttributeOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.ModifyInstanceAttributeOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ec2.ModifyInstanceAttributeInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSecurityGroup provides a mock function with given fields: _a0
func (_m *EC2er) DeleteSecurityGroup(_a0 *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
	ret := _m.Called(_a0)

	var r0 *ec2.DeleteSecurityGroupOutput
	if rf, ok := ret.Get(0).(func(*ec2.DeleteSecurityGroupInput) *ec2.DeleteSecurityGroupOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.DeleteSecurityGroupOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ec2.DeleteSecurityGroupInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

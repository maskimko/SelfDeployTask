package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	Subnet1Cidr string = "172.23.2.0/28"
	Subnet2Cidr string = "172.23.2.16/28"
)

func CreateSubnets(vpcId *string, svc *ec2.EC2) ([]*string, error) {
	subnets := make([]*string, 2, 2)
	azs, err := GetAzIds(svc)
	if err != nil {
		return subnets, err
	}

	subnets[0], err = createSubnet(aws.String(Subnet1Cidr), vpcId, azs[0], svc)
	if err != nil {
		return subnets, err
	}
	subnets[1], err = createSubnet(aws.String(Subnet2Cidr), vpcId, azs[1], svc)
	if err != nil {
		return subnets, err
	}
	return subnets, nil
}
func createSubnet(subnet, vpcId, az *string, svc *ec2.EC2) (*string, error) {

	input := &ec2.CreateSubnetInput{
		CidrBlock:        subnet,
		VpcId:            vpcId,
		AvailabilityZone: az}

	result, err := svc.CreateSubnet(input)
	subnetId := result.Subnet.SubnetId
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return nil, err
	}

	fmt.Println(result)
	err = LabelResource(subnetId, svc)
	if err != nil {
		log.Printf("Cannot label subnet %s: %s", *subnetId, err)
	}
	return subnetId, nil
}

func DeleteSubnets(subnetsIds []*string, svc *ec2.EC2) error {
	for _, net := range subnetsIds {
		err := DeleteSubnet(net, svc)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteSubnet(subnetId *string, svc *ec2.EC2) error {

	input := &ec2.DeleteSubnetInput{
		SubnetId: subnetId,
	}

	result, err := svc.DeleteSubnet(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return err
	}

	fmt.Println(result)
	return nil
}

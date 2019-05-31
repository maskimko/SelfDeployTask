package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	VpcCidrBlock string = "172.23.2.0/26"
)

func CreateVpc(vpcName string, ec2Service *ec2.EC2) (*string, error) {

	vpcOut, err := ec2Service.CreateVpc(&ec2.CreateVpcInput{CidrBlock: aws.String(VpcCidrBlock)})
	if err != nil {
		log.Printf("Cannot create VPC: %s", err)
		return nil, err
	}
	vpcId := vpcOut.Vpc.VpcId
	fmt.Printf("Creating VPC: id %s\n", *vpcId)
	descVpcIn := &ec2.DescribeVpcsInput{VpcIds: []*string{vpcId}}

	err = ec2Service.WaitUntilVpcExists(descVpcIn)
	if err != nil {
		log.Printf("VPC creation timeout: %s", err)
		return nil, err
	}
	fmt.Printf("VPC id: %s has been created successfully!\n", *vpcId)

	NameResource(&vpcName, vpcId, ec2Service)
	LabelResource(vpcId, ec2Service)
	//Describe to ensure that labes exist
	descVpcOut, err := ec2Service.DescribeVpcs(&ec2.DescribeVpcsInput{VpcIds: []*string{vpcId}})
	if err != nil {
		log.Printf("Cannot describe VPC: %s", err)
		return nil, err
	}
	fmt.Print(descVpcOut.String())

	return vpcId, nil
}

func GetVpcArn(region, accountId, vpcId string) string {
	vpcArnStr := fmt.Sprintf("arn:aws:ec2:%s:%s:vpc/%s", region, accountId, vpcId)
	fmt.Printf("VPC ARN: %s", vpcArnStr)
	return vpcArnStr
}

package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CreateVpc(vpcName string, ec2Service *ec2.EC2) error {

	vpcOut, err := ec2Service.CreateVpc(&ec2.CreateVpcInput{CidrBlock: aws.String("172.23.2.0/26")})
	if err != nil {
		log.Printf("Cannot create VPC: %s", err)
		return err
	}
	vpcId := vpcOut.Vpc.VpcId
	fmt.Printf("Creating VPC: id %s\n", *vpcId)
	descVpcIn := &ec2.DescribeVpcsInput{VpcIds: []*string{vpcId}}

	err = ec2Service.WaitUntilVpcExists(descVpcIn)
	if err != nil {
		log.Printf("Cannot describe VPC: %s", err)
		return err
	}
	fmt.Printf("VPC id: %s has been created successfully!\n", *vpcId)
	descVpcOut, err := ec2Service.DescribeVpcs(&ec2.DescribeVpcsInput{VpcIds: []*string{vpcId}})
	if err != nil {
		log.Printf("Cannot describe VPC: %s", err)
		return err
	}
	fmt.Print(descVpcOut.String())

	NameResource(&vpcName, vpcId, ec2Service)
	LabelResource(vpcId, ec2Service)
	descVpcOut, err = ec2Service.DescribeVpcs(&ec2.DescribeVpcsInput{VpcIds: []*string{vpcId}})
	if err != nil {
		log.Fatalf("Cannot describe VPC: %s", err)
	}
	fmt.Print(descVpcOut.String())
	return nil
}

func GetVpcArn(region, accountId, vpcId string) string {
	vpcArnStr := fmt.Sprintf("arn:aws:ec2:%s:%s:vpc/%s", region, accountId, vpcId)
	fmt.Printf("VPC ARN: %s", vpcArnStr)
	return vpcArnStr
}

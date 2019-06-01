package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CreateSecurityGroups(ipAddress, vpcId *string, svc *ec2.EC2) ([]*string, error) {
	sgIds := make([]*string, 2, 2)
	instanceSgId, err := CreateInstanceSecurityGroup(vpcId, svc)
	if err != nil {
		log.Printf("Cannot create security group for instance: %s", err)
		return sgIds, err
	}
	sgIds[0] = instanceSgId
	elbSgId, err := CreateELBSecurityGroup(vpcId, svc)
	if err != nil {
		log.Printf("Cannot create security group for balancer: %s", err)
		return sgIds, err
	}
	sgIds[1] = elbSgId
	err = AuthorizeInstanceAccess(instanceSgId, ipAddress, svc)
	if err != nil {
		log.Printf("Cannot create security group rules for instance: %s", err)
		return sgIds, err
	}
	err = AuthorizeElbAccess(elbSgId, ipAddress, svc)
	if err != nil {
		log.Printf("Cannot create security group rule for balancer: %s", err)
		return sgIds, err
	}
	return sgIds, nil
}

func CreateInstanceSecurityGroup(vpcId *string, svc *ec2.EC2) (*string, error) {

	input := &ec2.CreateSecurityGroupInput{
		Description: aws.String("Test instance security group"),
		GroupName:   aws.String("InstanceSecurityGroup"),
		VpcId:       vpcId,
	}

	result, err := svc.CreateSecurityGroup(input)
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
	return result.GroupId, nil
}

func CreateELBSecurityGroup(vpcId *string, svc *ec2.EC2) (*string, error) {
	input := &ec2.CreateSecurityGroupInput{
		Description: aws.String("Test balancer security group"),
		GroupName:   aws.String("ElbSecurityGroup"),
		VpcId:       vpcId,
	}

	result, err := svc.CreateSecurityGroup(input)
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
	return result.GroupId, nil
}

func AuthorizeInstanceAccess(sgId, ipAddress *string, svc *ec2.EC2) error {
	ipCidr := fmt.Sprintf("%s/32", *ipAddress)
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: sgId,
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(22),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(ipCidr),
						Description: aws.String("SSH access"),
					},
				},
				ToPort: aws.Int64(22),
			},
			{
				FromPort:   aws.Int64(1989),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(VpcCidrBlock),
						Description: aws.String("Access to the application TCP socket"),
					},
				},
				ToPort: aws.Int64(1989),
			},
		},
	}

	result, err := svc.AuthorizeSecurityGroupIngress(input)
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
		return err
	}

	fmt.Println(result)
	return nil
}

func AuthorizeElbAccess(sgId, ipAddress *string, svc *ec2.EC2) error {
	ipCidr := fmt.Sprintf("%s/32", *ipAddress)
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: sgId,
		IpPermissions: []*ec2.IpPermission{

			{
				FromPort:   aws.Int64(1989),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(ipCidr),
						Description: aws.String("Access to the application TCP socket"),
					},
				},
				ToPort: aws.Int64(1989),
			},
		},
	}

	result, err := svc.AuthorizeSecurityGroupIngress(input)
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
		return err
	}

	fmt.Println(result)
	return nil
}

func DeleteSecurityGroups(securityGroups []*string, svc *ec2.EC2) error {
	for _, sg := range securityGroups {
		err := DeleteSecurityGroup(sg, svc)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteSecurityGroup(sgId *string, svc *ec2.EC2) error {

	input := &ec2.DeleteSecurityGroupInput{
		GroupId: sgId,
	}

	result, err := svc.DeleteSecurityGroup(input)
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

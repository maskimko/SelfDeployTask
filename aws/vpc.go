package aws

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

func DeleteVpc(vpcId *string, svc *ec2.EC2) error {

	input := &ec2.DeleteVpcInput{
		VpcId: vpcId,
	}

	result, err := svc.DeleteVpc(input)
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

func CreateInternetGateway(svc *ec2.EC2) (*string, error) {

	input := &ec2.CreateInternetGatewayInput{}

	result, err := svc.CreateInternetGateway(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}
	igId := result.InternetGateway.InternetGatewayId
	fmt.Println(result)
	return igId, nil
}

func AttachIgw(igwId, vpcId *string, svc *ec2.EC2) error {

	input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: igwId,
		VpcId:             vpcId,
	}

	result, err := svc.AttachInternetGateway(input)
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

func RouteInternetTraffic(igwId, vpcId *string, svc *ec2.EC2) error {

	describeInput := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					vpcId,
				},
			},
		},
	}

	result, err := svc.DescribeRouteTables(describeInput)
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
	if len(result.RouteTables) == 0 {
		return errors.New("No route tables found")
	}
	if len(result.RouteTables) > 1 {
		log.Printf("Weird, there are more than 1 route table.\n\tI am oigng to pick first one from list and use it as a default one")
	}
	//Assuming that we have single default route table
	rtbId := result.RouteTables[0].RouteTableId
	err = NameResource(aws.String(ResourceName), rtbId, svc)
	if err != nil {
		log.Printf("Cannot name resource %s", *rtbId)
	}
	err = LabelResource(rtbId, svc)
	if err != nil {
		log.Printf("Cannot label resource %s", *rtbId)
	}

	routeInput := &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            igwId,
		RouteTableId:         rtbId,
	}

	routeResult, err := svc.CreateRoute(routeInput)
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

	fmt.Println(routeResult)
	return nil
}

func DetachInternetGateway(igwId, vpcId *string, timeoutSeconds int16, svc *ec2.EC2) error {

	describeInput := &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{igwId},
	}

	detachInput := &ec2.DetachInternetGatewayInput{
		InternetGatewayId: igwId,
		VpcId:             vpcId,
	}
	var sleep int16 = 5

	for {
		detachResult, err := svc.DetachInternetGateway(detachInput)
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

		fmt.Println(detachResult)
		describeResult, err := svc.DescribeInternetGateways(describeInput)
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

		fmt.Println(describeResult)
		if len(describeResult.InternetGateways) == 0 {
			log.Println("It looks there are no IGWs. It is weird but OK. Returning")
			break
		} else {
			igw := describeResult.InternetGateways[0]
			if len(igw.Attachments) == 0 {
				log.Println("Internet gateway is detached")
				break
			}
		}
		time.Sleep(time.Duration(sleep) * time.Second)
		sleep = sleep * 2
		if sleep > timeoutSeconds {
			return errors.New("Timeout waiting for Internet gateway to detach")

		}
	}
	return nil
}

func DeleteInternetGateway(igwId *string, svc *ec2.EC2) error {

	input := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: igwId,
	}

	result, err := svc.DeleteInternetGateway(input)
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

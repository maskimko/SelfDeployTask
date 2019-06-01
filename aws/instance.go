package aws

import (
	"fmt"
	"log"

	"wix/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func RunInstances(keyName *string, subnets []*string, sg *string, svc *ec2.EC2) ([]*string, error) {
	instanceIds := make([]*string, 2, 2)
	id, err := RunInstance(subnets[0], sg, keyName, svc)
	if err != nil {
		return nil, err
	}
	instanceIds[0] = id
	id, err = RunInstance(subnets[1], sg, keyName, svc)
	if err != nil {
		return nil, err
	}
	instanceIds[1] = id
	return instanceIds, nil
}

func RunInstance(subnet, sg, keyName *string, svc *ec2.EC2) (*string, error) {
	var instanceId *string
	region := svc.Config.Region
	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("xvdm"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(5),
				},
			},
		},
		ImageId:      aws.String(AMIMap[*region]),
		InstanceType: aws.String(DefaultInstanceType),
		KeyName:      keyName,
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		SecurityGroupIds: []*string{
			sg,
		},
		SubnetId: subnet,
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("TaskOwner"),
						Value: aws.String("maskimko"),
					},
				},
			},
		},
	}

	result, err := svc.RunInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}
	instanceId = result.Instances[0].InstanceId
	fmt.Println(result)
	fmt.Printf("Launched instance %s of type %s from AMI %s\n", *instanceId,
		*(result.Instances[0].InstanceType),
		*(result.Instances[0].ImageId))
	return instanceId, nil
}

func TerminateInstances(instances []*string, svc *ec2.EC2) error {

	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: instances}

	result, err := svc.TerminateInstances(terminateInput)
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

	waitInput := &ec2.DescribeInstancesInput{
		InstanceIds: instances}
	err = svc.WaitUntilInstanceTerminated(waitInput)
	if err != nil {
		log.Printf("Cannot wait for instance termination: %s", err)
		return err
	}
	return nil
}

func WaitForOkInstances(instanceIds []*string, svc *ec2.EC2) error {
	log.Printf("Waiting for instances %s to become ready...", utils.Slice2String(instanceIds))
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: instanceIds,
	}
	err := svc.WaitUntilInstanceStatusOk(input)
	return err
}

func WaitForInstances2Run(instanceIds []*string, svc *ec2.EC2) error {
	log.Printf("Waiting for instances %s to run...", utils.Slice2String(instanceIds))
	input := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	}
	err := svc.WaitUntilInstanceRunning(input)
	return err
}

func GetPublicIps(instanceIds []*string, svc *ec2.EC2) ([]*string, error) {
	ipAddresses := make([]*string, 0)
	err := WaitForOkInstances(instanceIds, svc)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	}

	result, err := svc.DescribeInstances(input)
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
		return ipAddresses, err
	}

	fmt.Println(result)
	for _, reserv := range result.Reservations {
		for _, instance := range reserv.Instances {
			ipAddresses = append(ipAddresses, instance.PublicIpAddress)
		}
	}
	return ipAddresses, nil
}

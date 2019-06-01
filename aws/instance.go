package aws

import (
	"fmt"
	"log"

	"github.com/fatih/color"

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
	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("xvdm"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(5),
				},
			},
		},
		ImageId:      aws.String(DefaultAMI),
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
	color.White("Launched instance %s of type %s from AMI %s", *instanceId,
		*(result.Instances[0].InstanceType),
		*(result.Instances[0].ImageId))
	return instanceId, nil
}

func TerminateInstance(instances []*string, svc *ec2.EC2) error {

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

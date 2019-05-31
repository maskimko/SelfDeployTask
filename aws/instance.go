package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func RunInstances(keyName *string, subnets []*string, sg *string, svc *ec2.EC2) error {
	err := RunInstance(subnets[0], sg, keyName, svc)
	if err != nil {
		return err
	}
	err = RunInstance(subnets[1], sg, keyName, svc)
	if err != nil {
		return err
	}
	return nil
}

func RunInstance(subnet, sg, keyName *string, svc *ec2.EC2) error {

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
						Key:   aws.String("Purpose"),
						Value: aws.String("test"),
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
		return err
	}

	fmt.Println(result)
	return nil
}

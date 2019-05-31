package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func RunInstances(svc *ec2.EC2) {

	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sdh"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(100),
				},
			},
		},
		ImageId:      aws.String(DefaultAMI),
		InstanceType: aws.String(DefaultInstanceType),
		KeyName:      aws.String("my-key-pair"),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		SecurityGroupIds: []*string{
			aws.String("sg-1a2b3c4d"),
		},
		SubnetId: aws.String("subnet-6e7f829e"),
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
		return
	}

	fmt.Println(result)
}

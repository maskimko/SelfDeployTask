package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const KeyPairName string = "TestTaskKeyPair"

func CreateKeyPair(svc *ec2.EC2) (*string, error) {

	input := &ec2.CreateKeyPairInput{
		KeyName: aws.String(KeyPairName),
	}

	result, err := svc.CreateKeyPair(input)
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

	fmt.Println(result)
	return result.KeyMaterial, nil
}

func DeleteKeyPair(svc *ec2.EC2) error {

	input := &ec2.DeleteKeyPairInput{
		KeyName: aws.String(KeyPairName),
	}

	result, err := svc.DeleteKeyPair(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	fmt.Println(result)
	return nil
}

package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
)

const (
	myNickName    string = "maskimko"
	tagIdentifier string = "TaskOwner"
)

func LabelResource(resourceId *string, svc *ec2.EC2) error {
	vpcTagIn := ec2.CreateTagsInput{Resources: []*string{resourceId},
		Tags: []*ec2.Tag{&ec2.Tag{Key: aws.String(tagIdentifier), Value: aws.String(myNickName)}}}
	vpcTagOut, err := svc.CreateTags(&vpcTagIn)
	if err != nil {
		log.Printf("Cannot tag resource: %s\n", err)
		return err
	}
	fmt.Println(vpcTagOut.String())
	return nil
}

func NameResource(name, resourceId *string, svc *ec2.EC2) error {
	vpcTagIn := ec2.CreateTagsInput{Resources: []*string{resourceId},
		Tags: []*ec2.Tag{&ec2.Tag{Key: aws.String("Name"), Value: name}}}
	vpcTagOut, err := svc.CreateTags(&vpcTagIn)
	if err != nil {
		log.Printf("Cannot tag resource: %s\n", err)
		return err
	}
	fmt.Println(vpcTagOut.String())
	return nil
}

func LabelElb(name *string, svc *elb.ELB) error {

	input := &elb.AddTagsInput{
		LoadBalancerNames: []*string{
			aws.String(ElbName),
		},
		Tags: []*elb.Tag{
			{
				Key:   aws.String("Name"),
				Value: name,
			},
			{
				Key:   aws.String(tagIdentifier),
				Value: aws.String(myNickName),
			},
		},
	}

	result, err := svc.AddTags(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeAccessPointNotFoundException:
				log.Println(elb.ErrCodeAccessPointNotFoundException, aerr.Error())
			case elb.ErrCodeTooManyTagsException:
				log.Println(elb.ErrCodeTooManyTagsException, aerr.Error())
			case elb.ErrCodeDuplicateTagKeysException:
				log.Println(elb.ErrCodeDuplicateTagKeysException, aerr.Error())
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

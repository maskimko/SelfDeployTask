package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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

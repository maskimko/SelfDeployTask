package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func BindPublicIps(instanceIds []*string, svc *ec2.EC2) ([]*string, []*string, error) {
	alls := make([]*string, 0)
	asss := make([]*string, 0)
	for _, i := range instanceIds {
		all, ass, err := BindPublicIp(i, svc)
		if err != nil {
			return nil, nil, err
		}
		alls = append(alls, all)
		asss = append(asss, ass)
	}
	return alls, asss, nil
}

func BindPublicIp(instanceId *string, svc *ec2.EC2) (*string, *string, error) {
	var allocationId *string
	allocationInput := &ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	}

	allocationResult, err := svc.AllocateAddress(allocationInput)
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
		return nil, nil, err
	}

	fmt.Println(allocationResult)
	allocationId = allocationResult.AllocationId

	associationInput := &ec2.AssociateAddressInput{
		AllocationId: allocationId,
		InstanceId:   instanceId,
	}
	err = WaitForInstances2Run([]*string{instanceId}, svc)
	if err != nil {
		return allocationId, nil, err
	}
	associationResult, err := svc.AssociateAddress(associationInput)
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
		return allocationId, nil, err
	}
	associationId := associationResult.AssociationId
	fmt.Println(associationResult)
	return allocationId, associationId, nil
}

func UnbindPublicIps(allocationIds, associationIds []*string, svc *ec2.EC2) error {
	for _, ass := range associationIds {
		err := DisassociatePublicIp(ass, svc)
		if err != nil {
			return err
		}
	}
	for _, all := range allocationIds {
		err := ReleasePublicIp(all, svc)
		if err != nil {
			return err
		}
	}
	return nil
}

func DisassociatePublicIp(associationId *string, svc *ec2.EC2) error {
	disassociateInput := &ec2.DisassociateAddressInput{
		AssociationId: associationId,
	}

	result, err := svc.DisassociateAddress(disassociateInput)
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

func ReleasePublicIp(allocationId *string, svc *ec2.EC2) error {

	releaseInput := &ec2.ReleaseAddressInput{
		AllocationId: allocationId,
	}

	result, err := svc.ReleaseAddress(releaseInput)
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

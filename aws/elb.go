package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
)

const ElbName string = "TestELB"

func CreateElb(subnets []*string, sg *string, svc *elb.ELB) (*string, error) {

	input := &elb.CreateLoadBalancerInput{
		Listeners: []*elb.Listener{
			{
				InstancePort:     aws.Int64(1989),
				InstanceProtocol: aws.String("TCP"),
				LoadBalancerPort: aws.Int64(1989),
				Protocol:         aws.String("TCP"),
			},
		},
		LoadBalancerName: aws.String(ElbName),
		SecurityGroups: []*string{
			sg,
		},
		Subnets: subnets,
	}

	result, err := svc.CreateLoadBalancer(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeDuplicateAccessPointNameException:
				log.Println(elb.ErrCodeDuplicateAccessPointNameException, aerr.Error())
			case elb.ErrCodeTooManyAccessPointsException:
				log.Println(elb.ErrCodeTooManyAccessPointsException, aerr.Error())
			case elb.ErrCodeCertificateNotFoundException:
				log.Println(elb.ErrCodeCertificateNotFoundException, aerr.Error())
			case elb.ErrCodeInvalidConfigurationRequestException:
				log.Println(elb.ErrCodeInvalidConfigurationRequestException, aerr.Error())
			case elb.ErrCodeSubnetNotFoundException:
				log.Println(elb.ErrCodeSubnetNotFoundException, aerr.Error())
			case elb.ErrCodeInvalidSubnetException:
				log.Println(elb.ErrCodeInvalidSubnetException, aerr.Error())
			case elb.ErrCodeInvalidSecurityGroupException:
				log.Println(elb.ErrCodeInvalidSecurityGroupException, aerr.Error())
			case elb.ErrCodeInvalidSchemeException:
				log.Println(elb.ErrCodeInvalidSchemeException, aerr.Error())
			case elb.ErrCodeTooManyTagsException:
				log.Println(elb.ErrCodeTooManyTagsException, aerr.Error())
			case elb.ErrCodeDuplicateTagKeysException:
				log.Println(elb.ErrCodeDuplicateTagKeysException, aerr.Error())
			case elb.ErrCodeUnsupportedProtocolException:
				log.Println(elb.ErrCodeUnsupportedProtocolException, aerr.Error())
			case elb.ErrCodeOperationNotPermittedException:
				log.Println(elb.ErrCodeOperationNotPermittedException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return nil, err
	}

	fmt.Println(result)
	err = LabelElb(aws.String(ElbName), svc)
	if err != nil {
		log.Printf("Cannot label ELB: %s", err)
		return result.DNSName, err
	}
	return result.DNSName, nil
}

func AttachInstances2Elb(instances []*string, svc *elb.ELB) error {

	vms := make([]*elb.Instance, len(instances))
	for i, vm := range instances {
		vms[i] = &elb.Instance{InstanceId: vm}
	}
	input := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances:        vms,
		LoadBalancerName: aws.String(ElbName),
	}

	result, err := svc.RegisterInstancesWithLoadBalancer(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeAccessPointNotFoundException:
				log.Println(elb.ErrCodeAccessPointNotFoundException, aerr.Error())
			case elb.ErrCodeInvalidEndPointException:
				log.Println(elb.ErrCodeInvalidEndPointException, aerr.Error())
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

func DeleteElb(svc *elb.ELB) error {

	input := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(ElbName),
	}

	result, err := svc.DeleteLoadBalancer(input)
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

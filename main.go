package main
import (
	"log"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/aws/arn"
)
const VpcName string = "Go4learn"

func main(){
    sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	conf := defaults.Config()
	region := *(conf.Region)
	fmt.Println(region)
	iamService := iam.New(sess)
	userOut, err := iamService.GetUser(&iam.GetUserInput{})
	if err != nil {
		log.Fatal("Cannot create VPC: %s", err)
	}
	userArn, err := arn.Parse(*(userOut.User.Arn))
	if err != nil {
		log.Fatal("Cannot create VPC: %s", err)
	}
	accountId := userArn.AccountID
	fmt.Println(accountId)
	ec2Service := ec2.New(sess)
	vpcOut, err := ec2Service.CreateVpc(&ec2.CreateVpcInput{CidrBlock : aws.String("172.23.2.0/26")})
	if err != nil {
		log.Fatal("Cannot create VPC: %s", err)
	}
	vpcId := vpcOut.Vpc.VpcId
	fmt.Printf("Creating VPC: id %s\n", *vpcId)
	descVpcIn := &ec2.DescribeVpcsInput{VpcIds: []*string{vpcId}}
	
	err = ec2Service.WaitUntilVpcExists(descVpcIn)
	if err != nil {
		log.Fatal("Cannot describe VPC: %s", err)
	}
	fmt.Printf("VPC id: %s has been created successfully!\n", *vpcId)
	descVpcOut, err := ec2Service.DescribeVpcs(&ec2.DescribeVpcsInput{VpcIds: []*string{vpcId}})
	if err != nil {
		log.Fatal("Cannot describe VPC: %s", err)
	}
	fmt.Print(descVpcOut.String())
	
	vpcArnStr := fmt.Sprintf("arn:aws:ec2:%s:%s:vpc/%s",region,accountId,*vpcId)
	fmt.Printf("VPC ARN: %s", vpcArnStr)
	vpcTagIn  := ec2.CreateTagsInput{Resources: []*string{vpcId}, 
	Tags: []*ec2.Tag{&ec2.Tag{Key: aws.String("Name"), Value: aws.String(VpcName)}}}
	vpcTagOut, err := ec2Service.CreateTags(&vpcTagIn)
	if err != nil {
		log.Fatal("Cannot describe VPC: %s", err)
	}
	fmt.Print(vpcTagOut.String())
	descVpcOut, err = ec2Service.DescribeVpcs(&ec2.DescribeVpcsInput{VpcIds: []*string{vpcId}})
	if err != nil {
		log.Fatal("Cannot describe VPC: %s", err)
	}
	fmt.Print(descVpcOut.String())
}
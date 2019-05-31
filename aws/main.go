package aws

import (
	"fmt"
	"log"
	"wix/ssh"
	"wix/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

//CentOS 7 (x86_64) - with Updates HVM
const (
	DefaultAMI          string = "ami-02eac2c0129f6376b"
	DefaultInstanceType string = "t2.nano"
)

func GetDefaultSession() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}

func GetSession(region *string) *session.Session {
	if region == nil {
		return GetDefaultSession()
	} else {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: region}))
		return sess
	}
}

func Init(sshConf *ssh.SshConfig) error {
	//Get session
	sess := GetDefaultSession()
	conf := defaults.Config()
	region := *(conf.Region)

	iamService := iam.New(sess)
	userOut, err := iamService.GetUser(&iam.GetUserInput{})
	if err != nil {
		log.Printf("Cannot get current user info: %s\n", err)
		return err
	}
	userArn, err := arn.Parse(*(userOut.User.Arn))
	if err != nil {
		log.Printf("Cannot get current user ARN: %s\n", err)
		return err
	}
	accountId := userArn.AccountID
	fmt.Printf("Using AWS account id: %s\n", accountId)
	ec2Service := ec2.New(sess)
	azs, err := GetAzIds(ec2Service)
	if err != nil {
		return err
	}
	fmt.Printf("Using region %s. This region has these availability zones: %v\n", region, azs)
	vpcId, err := CreateVpc("Test-Vpc", ec2Service)
	if err != nil {
		return err
	}
	privKey, err := CreateKeyPair(ec2Service)
	sshConf.PrivateKey = []byte(*privKey)
	if err != nil {
		return err
	}
	subnets, err := CreateSubnets(vpcId, ec2Service)
	if err != nil {
		return err
	}
	fmt.Printf("Created subnets %v", subnets)
	ipAddr, err := utils.GetMyIp()
	if err != nil {
		return err
	}
	securityGroupIds, err := CreateSecurityGroups(ipAddr, vpcId, ec2Service)
	if err != nil {
		return err
	}
	fmt.Printf("Created security groups %v", securityGroupIds)
	err = RunInstances(aws.String(KeyPairName), subnets, securityGroupIds[0], ec2Service)
	if err != nil {
		return err
	}
	return nil
}

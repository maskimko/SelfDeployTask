package aws

import (
	"log"
	"wix/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/fatih/color"
)

//CentOS 7 (x86_64) - with Updates HVM
const (
	DefaultAMI          string = "ami-02eac2c0129f6376b"
	DefaultInstanceType string = "t2.nano"
)

type Inventory struct {
	VpcId          *string
	Instances      []*string
	PrivateKey     *string
	Session        *session.Session
	Region         *string
	IgwId          *string
	SecurityGroups []*string
	Subnets        []*string
}

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

func Init(inventory *Inventory) error {
	sess := inventory.Session
	conf := sess.Config
	region := conf.Region
	inventory.Region = region
	color.Blue("Using region '%s'", *region)
	iamService := iam.New(sess)
	ec2Service := ec2.New(sess)
	elbService := elb.New(sess)
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
	color.Green("Using AWS account id: '%s'", accountId)

	azs, err := GetAzIds(ec2Service)
	if err != nil {
		return err
	}
	color.Green("Using region '%s'. This region has these availability zones: %s", *region, utils.Slice2String(azs))
	vpcId, err := CreateVpc("Test-Vpc", ec2Service)
	if err != nil {
		return err
	}
	inventory.VpcId = vpcId
	color.Green("VPC '%s' has been successfully created", *vpcId)
	igwId, err := CreateInternetGateway(ec2Service)
	if err != nil {
		return err
	}
	inventory.IgwId = igwId
	color.Green("Internet Gateway '%s' has been successfully created", *igwId)
	AttachIgw(igwId, vpcId, ec2Service)
	if err != nil {
		return err
	}
	color.Green("Internet Gateway '%s' has been successfully attached to the VPC '%s'", *igwId, *vpcId)
	privKey, err := CreateKeyPair(ec2Service)
	inventory.PrivateKey = privKey
	if err != nil {
		return err
	}
	color.Green("SSH key pair '%s' has been successfully created", KeyPairName)
	subnets, err := CreateSubnets(vpcId, ec2Service)
	if err != nil {
		return err
	}
	inventory.Subnets = subnets
	color.Green("Created subnets %s", utils.Slice2String(subnets))
	ipAddr, err := utils.GetMyIp()
	if err != nil {
		return err
	}

	securityGroupIds, err := CreateSecurityGroups(ipAddr, vpcId, ec2Service)
	if err != nil {
		return err
	}
	inventory.SecurityGroups = securityGroupIds
	color.Green("Created security groups %s", utils.Slice2String(securityGroupIds))
	instanceIds, err := RunInstances(aws.String(KeyPairName), subnets, securityGroupIds[0], ec2Service)
	if err != nil {
		return err
	}
	inventory.Instances = instanceIds
	color.Green("Created instances '%v'", utils.Slice2String(instanceIds))
	dns, err := CreateElb(subnets, securityGroupIds[1], elbService)
	if err != nil {
		return err
	}
	color.Green("Created ELB. It is available by tcp://%s:1989", *dns)
	return nil
}

func Destroy(inventory *Inventory) error {
	color.Red("Destroying region %s...", *(inventory.Region))
	sess := inventory.Session
	ec2Service := ec2.New(sess)
	elbService := elb.New(sess)
	color.Red("Terminating instances %s. Please wait...", utils.Slice2String(inventory.Instances))
	err := TerminateInstance(inventory.Instances, ec2Service)
	if err != nil {
		return err
	}
	color.Red("Instances %s have been successfully terminated", utils.Slice2String(inventory.Instances))
	err = DeleteElb(elbService)
	if err != nil {
		return err
	}

	color.Red("ELB %s have been successfully removed", ElbName)
	err = DeleteSecurityGroups(inventory.SecurityGroups, ec2Service)
	if err != nil {
		return err
	}
	color.Red("Security groups %s have been successfully removed", utils.Slice2String(inventory.SecurityGroups))
	err = DeleteSubnets(inventory.Subnets, ec2Service)
	if err != nil {
		return err
	}
	color.Red("Subnets %s have been successfully removed", utils.Slice2String(inventory.Subnets))

	err = DetachInternetGateway(inventory.IgwId, inventory.VpcId, ec2Service)
	if err != nil {
		return err
	}
	err = WaitUntilIgwDetached(inventory.IgwId, 120, ec2Service)
	if err != nil {
		return err
	}
	color.Red("Internet gateway %s have been successfully detached from VPC %s", *(inventory.IgwId), *(inventory.VpcId))
	DeleteInternetGateway(inventory.IgwId, ec2Service)
	if err != nil {
		return err
	}
	color.Red("IGW %s have been successfully removed", *(inventory.IgwId))
	err = DeleteVpc(inventory.VpcId, ec2Service)
	if err != nil {
		return err
	}
	color.Red("VPC %s have been successfully removed", inventory.VpcId)
	color.Magenta("Region %s is clean", *(inventory.Region))
	return nil
}

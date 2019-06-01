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
	// DefaultAMI          string = "ami-02eac2c0129f6376b"
	DefaultInstanceType string = "t2.nano"
)

var AMIMap map[string]string = map[string]string{
	"ap-northeast-1": "ami-25bd2743",
	"ap-northeast-2": "ami-7248e81c",
	"ap-south-1":     "ami-5d99ce32",
	"ap-southeast-1": "ami-d2fa88ae",
	"ap-southeast-2": "ami-b6bb47d4",
	"ca-central-1":   "ami-dcad28b8",
	"eu-central-1":   "ami-337be65c",
	"eu-west-1":      "ami-6e28b517",
	"eu-west-2":      "ami-ee6a718a",
	"eu-west-3":      "ami-bfff49c2",
	"sa-east-1":      "ami-f9adef95",
	"us-east-1":      "ami-4bf3d731",
	"us-east-2":      "ami-e1496384",
	"us-west-1":      "ami-65e0e305",
	"us-west-2":      "ami-a042f4d8",
}

type Inventory struct {
	VpcId          *string
	Instances      []*string
	PublicIps      []*string
	AllocationIds  []*string
	AssociationIds []*string
	PrivateKey     *string
	Session        *session.Session
	Region         *string
	IgwId          *string
	SecurityGroups []*string
	Subnets        []*string
}

func (i *Inventory) Clone() *Inventory {
	inv := Inventory{
		VpcId:          i.VpcId,
		Instances:      i.Instances,
		PublicIps:      i.PublicIps,
		PrivateKey:     i.PrivateKey,
		AllocationIds:  i.AllocationIds,
		AssociationIds: i.AssociationIds,
		Session:        i.Session,
		Region:         i.Region,
		IgwId:          i.IgwId,
		SecurityGroups: i.SecurityGroups,
		Subnets:        i.Subnets,
	}
	return &inv
}

func (i *Inventory) GetPrivateKey() *[]byte {
	pKey := []byte(*i.PrivateKey)
	return &pKey
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
  color.HiBlue("For debugging purpose you may want to save this private key:\n%s\n", *privKey)
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
	allocationIds, associationIds, err := BindPublicIps(instanceIds, ec2Service)
	if err != nil {
		return err
	}
	inventory.AllocationIds = allocationIds
	inventory.AssociationIds = associationIds
	publicIps, err := GetPublicIps(instanceIds, ec2Service)
	if err != nil {
		return err
	}
	inventory.PublicIps = publicIps
	color.Green("Created instances '%s' with public IP addresses (%s)",
		utils.Slice2String(instanceIds),
		utils.Slice2String(publicIps))
	dns, err := CreateElb(subnets, securityGroupIds[1], elbService)
	if err != nil {
		return err
	}
	color.Green("Created ELB. It is available by tcp://%s:1989", *dns)

	color.Cyan("Region %s has been successfully initialized!", *region)
	return nil
}

func Destroy(inventory *Inventory) error {
	color.Red("Destroying region %s...", *(inventory.Region))
	sess := inventory.Session
	ec2Service := ec2.New(sess)
	elbService := elb.New(sess)

	err := UnbindPublicIps(inventory.AllocationIds, inventory.AssociationIds, ec2Service)
	if err != nil {
		return err
	}
	color.Red("Elastic IP addresses %s\n\t(allocation ids %s,\n\tassociation ids %s)\n\thas been successfully released",
		utils.Slice2String(inventory.PublicIps),
		utils.Slice2String(inventory.AllocationIds),
    utils.Slice2String(inventory.AssociationIds))
	color.Red("Terminating instances %s. Please wait...", utils.Slice2String(inventory.Instances))

	err = TerminateInstances(inventory.Instances, ec2Service)
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

	err = DetachInternetGateway(inventory.IgwId, inventory.VpcId, 120, ec2Service)
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
	err = DeleteKeyPair(ec2Service)
	if err != nil {
		return err
	}
	color.Red("Key pair %s have been successfully removed", KeyPairName)
	color.Magenta("Region %s is clean", *(inventory.Region))
	return nil
}

package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
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

func Init() error {
	//Get session
	sess := GetDefaultSession()
	conf := defaults.Config()
	region := *(conf.Region)
	fmt.Println(region)
	iamService := iam.New(sess)
	userOut, err := iamService.GetUser(&iam.GetUserInput{})
	if err != nil {
		log.Printf("Cannot get current user info: %s", err)
		return err
	}
	userArn, err := arn.Parse(*(userOut.User.Arn))
	if err != nil {
		log.Printf("Cannot get current user ARN: %s", err)
		return err
	}
	accountId := userArn.AccountID
	fmt.Println(accountId)
	return nil
}

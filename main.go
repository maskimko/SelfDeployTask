package main

import (
	"log"
	"wix/aws"
	"wix/server"
	"wix/ssh"
	"wix/utils"

	"github.com/fatih/color"
)

func main() {
	color.White("Welcome to the wix task of Maksym Shkolnyi")
	me, err := utils.GetPath2Itself()
	if err != nil {
		log.Fatalf("Cannot get path to itself: %s\n", err)
	}
	color.Yellow("This executable is located at: %s\n", me)
	myIp, err := utils.GetMyIp()
	if err != nil {
		log.Fatalf("Cannot get my IP address: %s\n", err)
	}

	color.White("My external IP address is: %s\n", *myIp)
	// ssh.CopyItself(sshConfig)
	// fmt.Println("I moved myself to the remote machine")
	// ssh.RunCommand("uptime", sshConfig)
	// fmt.Println("I launched myself to the remote machine")

	color.HiBlue("Checking where I am. Please wait...")
	inAWS := utils.AmIinAnAWS()
	if inAWS {
		color.Yellow("I exist in AWS")
	} else {
		color.HiBlue("I defenitely exist not in AWS. Perhaps on your laptop;)")
	}
	awsSession := aws.GetDefaultSession()
	awsInventory := &aws.Inventory{Session: awsSession,
		InAWS: inAWS,
	}
	err = aws.Init(awsInventory)
	if err != nil {
		log.Fatalf("Cannot perform AWS initialization: %s", err)
	}

	//Deploy itself
	err = ssh.Deploy(awsInventory.PublicIps, awsInventory.GetPrivateKey())
	if err != nil {
		// log.Fatalf("Cannot deploy myself to servers: %s", err)
		log.Printf("Cannot deploy myself to servers: %s", err)
	}
	color.Green("Application has been successfully deployed to the servers\n\tAvailable endpoints:")
	for _, ip := range awsInventory.PublicIps {
		color.Cyan("\t%s:1989", *ip)
	}

	//String server
	err = server.Start(1989, awsInventory)
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}

}

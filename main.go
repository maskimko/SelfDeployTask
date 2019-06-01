package main

import (
	"fmt"
	"log"
	"wix/aws"
	"wix/server"
	"wix/utils"
)

const VpcName string = "Go4learn"

func main() {
	me, err := utils.GetPath2Itself()
	if err != nil {
		log.Fatalf("Cannot get path to itself: %s\n", err)
	}
	fmt.Printf("This executable is located at: %s\n", me)
	myIp, err := utils.GetMyIp()
	if err != nil {
		log.Fatalf("Cannot get my IP address: %s\n", err)
	}

	// sshConfig := &ssh.SshConfig{Host: "ansible.tonicfordev.com",
	// 	Port:    22,
	// 	KeyPath: "/Users/maksym.shkolnyi/.ssh/tonic",
	// 	User:    "maksym.shkolnyi"}
	//
	fmt.Printf("My external IP address is: %s\n", *myIp)
	// ssh.CopyItself(sshConfig)
	// fmt.Println("I moved myself to the remote machine")
	// ssh.RunCommand("uptime", sshConfig)
	// fmt.Println("I launched myself to the remote machine")

	fmt.Println("Checking where I am. Please wait...")
	inAWS := utils.AmIinAnAWS()
	if inAWS {
		fmt.Println("I exist in AWS")
	} else {
		fmt.Println("I defenitely exist not in AWS. Perhaps on your laptop;)")
	}
	//sshConfig := &ssh.SshConfig{}

	awsSession := aws.GetDefaultSession()
	awsInventory := &aws.Inventory{Session: awsSession}
	err = aws.Init(awsInventory)
	if err != nil {
		log.Fatalf("Cannot perform AWS initialization: %s", err)
	}
	err = server.Start(1989, awsInventory)
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
}

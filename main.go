package main

import (
	"fmt"
	"log"
	"wix/server"
	"wix/utils"
)

const VpcName string = "Go4learn"

func main() {
	me, err := utils.GetPath2Itself()
	if err != nil {
		log.Fatalf("Cannot get path to itself: %s", err)
	}
	fmt.Printf("This executable is located at: %s\n", me)
	myIp, err := utils.GetMyIp()
	if err != nil {
		log.Fatalf("Cannot get my IP address: %s", err)
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
	err = server.Start(1989)
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
}

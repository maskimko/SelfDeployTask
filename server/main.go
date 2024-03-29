package server

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"wix/aws"
	"wix/ssh"
)

func Start(port int16, inventory *aws.Inventory) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Cannot bind port %d: %s\n", port, err)
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Connection error: %s\n", err)
		}
		go handleConnection(conn, inventory)
	}
}

func handleConnection(conn net.Conn, inventory *aws.Inventory) error {
	defer conn.Close()
	rAddr := conn.RemoteAddr()
	fmt.Printf("get Connection from %s\n", rAddr)
	b, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Printf("Cannot read message: %s\n", err)
		return err
	}
	log.Printf("Received message %s (length %d)\n", b, len(b))
	err = dispatch(&b, inventory)
	if err != nil {
		log.Printf("Got error while handling message: %s\n", err)
	}
	return err
}

func dispatch(message *[]byte, inventory *aws.Inventory) error {
	stopRegex := regexp.MustCompile("^stop$")
	moveRegex := regexp.MustCompile("^moveto '([a-z0-9-]+)'$")
	rows := bytes.Split(*message, []byte("\n"))
	for rn, row := range rows {
		if len(row) > 0 {
			if stopRegex.Match(row) {
				err := handleStop(inventory)
				if err != nil {
					return err
				}
				continue
			}
			if moveRegex.Match(row) {
				matches := moveRegex.FindSubmatch(row)
				if len(matches) > 1 {
					err := handleMove(matches[1], inventory)
					if err != nil {
						return err
					}
				} else {
					log.Printf("Cannot extract region from the move command %s\n", row)
				}
				continue
			}
			return errors.New("Unsupported command")
		} else {
			log.Printf("Skipping an empty row %d", rn)
		}
	}
	return nil
}

func handleStop(inventory *aws.Inventory) error {
	err := aws.Destroy(true, inventory)
	return err
}

func handleStopAfterMove(inventory *aws.Inventory) error {
	err := aws.Destroy(false, inventory)
	return err
}

func handleMove(r []byte, inventory *aws.Inventory) error {
	//Cannot avoid type conversion here, as it requires to change public interface
	// and do the type conversion inside the GetSession function
	region := string(r)
	log.Printf("Received move to region %s signal. Start initializing a new region %s\n", region, region)
	oldInventory := inventory.Clone()
	awsSession := aws.GetSession(&region)
	inventory.Session = awsSession
	err := aws.Init(inventory)
	if err != nil {
		return err
	}
	err = ssh.Deploy(inventory.PublicIps, inventory.GetPrivateKey())
	if err != nil {
		log.Printf("Cannot deploy myself to servers: %s", err)
	}
	log.Printf("Shutting down deployment in initial region %s", *(oldInventory.Region))
	err = handleStopAfterMove(oldInventory)
	if err != nil {
		return err
	}
	return nil
}

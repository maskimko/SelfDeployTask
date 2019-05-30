package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

func Start(port int16) error {
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()
	rAddr := conn.RemoteAddr()
	fmt.Printf("get Connection from %s\n", rAddr)
	messageBuf := make([]byte, 0)
	tmp := make([]byte, 256)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Printf("Cannot read message: %s\n", err)
				return err
			}
			break
		}
		messageBuf = append(messageBuf, tmp[:n]...)
	}
	message := string(messageBuf[:len(messageBuf)])
	log.Printf("Received message %s (length %d)\n", message, len(message))
	err := dispatch(&message)
	if err != nil {
		log.Printf("Got error while handling message: %s\n", err)
	}
	return err
}

func dispatch(message *string) error {
	stopRegex, _ := regexp.Compile("^stop$")
	moveRegex, _ := regexp.Compile("^moveto '([a-z0-9-]+)'$")
	rows := strings.Split(*message, "\n")
	for rn, row := range rows {
		if len(row) > 0 {
			if stopRegex.MatchString(row) {
				err := handleStop()
				if err != nil {
					return err
				}
				continue
			}
			if moveRegex.MatchString(row) {
				matches := moveRegex.FindStringSubmatch(row)
				if len(matches) > 1 {
					err := handleMove(matches[1])
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

func handleStop() error {
	fmt.Println("Received stop signal (NOT implemented)")
	return nil
}

func handleMove(region string) error {
	fmt.Printf("Received move to region %s signal (NOT implemented)\n", region)
	return nil
}

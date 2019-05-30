package ssh

import (
	"fmt"
	"os"
	"wix/utils"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func CopyItself(configuration *SshConfig) error {
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PrivateKey(configuration.User, configuration.KeyPath, ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(configuration.GetHostPort(), &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establisch a connection to the remote server ", err)
		return err
	}

	// Open a file
	me, err := utils.GetPath2Itself()
	f, _ := os.Open(me)

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)
	remotePath := utils.GetRemotePath(configuration.User, me)
	err = client.CopyFile(f, remotePath, "0755")

	if err != nil {
		fmt.Println("Error while copying file ", err)
		return err
	}
	return nil
}

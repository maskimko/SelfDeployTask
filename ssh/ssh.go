package ssh

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func PublicKeyFromFile(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Cannot read private ssh key: %s", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Cannot parse private ssh key: %s", err)
	}
	return ssh.PublicKeys(signer)
}

//An unencrypted PEM encoded RSA private key.
func publicKey(privateKey *[]byte) ssh.AuthMethod {

	signer, err := ssh.ParsePrivateKey(*privateKey)
	if err != nil {
		log.Fatalf("Cannot parse private ssh key: %s", err)
	}
	return ssh.PublicKeys(signer)
}

func RunCommand(cmd string, configuration *SshConfig) error {

	config := &ssh.ClientConfig{
		User: configuration.User,
		Auth: []ssh.AuthMethod{
			publicKey(&configuration.PrivateKey)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", configuration.GetHostPort(), config)

	if err != nil {
		log.Printf("Cannot establish ssh connection: %s", err)
		return err
	}
	defer conn.Close()

	sess, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, sessStderr)

	err = sess.Run(cmd)
	if err != nil {
		return err
	}
	return nil
}

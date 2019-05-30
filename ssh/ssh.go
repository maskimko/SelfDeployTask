package ssh

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func publicKey(path string) ssh.AuthMethod {
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

func RunCommand(cmd string, configuration *SshConfig) {

	config := &ssh.ClientConfig{
		User: configuration.User,
		Auth: []ssh.AuthMethod{
			publicKey(configuration.KeyPath)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", configuration.GetHostPort(), config)

	if err != nil {
		log.Fatalf("Cannot establish ssh connection: %s", err)
	}
	defer conn.Close()

	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, sessStderr)

	err = sess.Run(cmd)
	if err != nil {
		panic(err)
	}
}

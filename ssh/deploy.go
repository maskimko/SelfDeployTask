package ssh

import (
	"fmt"
	"log"
	"wix/utils"
)

const (
	DefaultSshUser string = "centos"
	Executable     string = "wix-task"
	unitTemplate   string = `
[Unit]
Description=Test task daemon (author Maksym Shkolnyi)


[Service]
Type=forking
PIDFile=/var/run/%s.pid
#EnvironmentFile=-/etc/somefile
ExecStart=/usr/local/bin/%s --pid-file=/var/run/%s

[Install]
WantedBy=multi-user.target
`
)

func Deploy(ips []*string, privateKey *[]byte) error {
	for _, ip := range ips {
		sshConfig := &SshConfig{
			Host:       *ip,
			Port:       22,
			PrivateKey: *privateKey,
			User:       "centos"}
		err := DeployHost(sshConfig)
		if err != nil {
			log.Printf("Cannot deploy host %s", *ip)
			return err
		}
	}
	return nil
}

func DeployHost(sc *SshConfig) error {
	err := CopyItself(sc)
	if err != nil {
		return err
	}
	me, err := utils.GetPath2Itself()
	if err != nil {
		return err
	}
	remote := utils.GetRemotePath(sc.User, me)

	err = RunCommand(fmt.Sprintf("sudo install -m 0755 %s /usr/local/bin/%s", remote, Executable), sc)
	if err != nil {
		log.Printf("Cannot install binary: %s", err)
		return err
	}
	unitFileContents := fmt.Sprintf(unitTemplate, Executable, Executable, Executable)
	err = RunCommand(fmt.Sprintf("cat '%s' > /tmp/%s.unitfile", unitFileContents, Executable), sc)
	if err != nil {
		log.Printf("Cannot create temp unit file: %s", err)
		return err
	}
	err = RunCommand(fmt.Sprintf("sudo install -m 0644 /tmp/%s.unitfile /use/lib/systemd/system/%s", Executable, Executable), sc)
	if err != nil {
		log.Printf("Cannot install unit file: %s", err)
		return err
	}
	err = RunCommand("sudo systemctl daemon-reload", sc)
	if err != nil {
		log.Printf("Cannnot reload systemd: %s", err)
		return err
	}
	err = RunCommand("sudo systemctl start ", sc)
	if err != nil {
		log.Printf("Cannot install unit file: %s", err)
		return err
	}
	return nil
}

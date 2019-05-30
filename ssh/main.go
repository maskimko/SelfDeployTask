package ssh

import (
	"fmt"
)

type SshConfig struct {
	Host    string
	Port    int16
	User    string
	KeyPath string
}

func (sc *SshConfig) GetHostPort() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}

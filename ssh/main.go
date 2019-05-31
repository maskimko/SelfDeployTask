package ssh

import (
	"fmt"
	"io/ioutil"
	"log"
)

type SshConfig struct {
	Host string
	Port int16
	User string
	//An unencrypted PEM encoded RSA private key.
	PrivateKey []byte
}

func (sc *SshConfig) GetHostPort() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}

func (sc *SshConfig) WithPrivateKeyFromFile(path string) (*SshConfig, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Cannot read private ssh key: %s", err)
		return nil, err
	}
	sc.PrivateKey = key
	return sc, nil
}

package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

const MagicIp string = "http://169.254.169.254/latest/meta-data/instance-id"

func AmIinAnAWS() bool {

	timeout := time.Second * 3
	var client = &http.Client{
		Timeout: timeout}
	resp, err := client.Get(MagicIp)
	if err != nil {
		log.Printf("Cannot reach out to %s", MagicIp)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Cannot read response from %s", MagicIp)
		return false
	}
	instanceId := string(body)
	instanceRegex, _ := regexp.Compile("^i-[0-9a-f]{8,}$")
	return instanceRegex.MatchString(instanceId)
}

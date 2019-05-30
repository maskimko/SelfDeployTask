package utils

import (
	"io/ioutil"
	"log"
	"net/http"
)

const AwsIpCheckerUrl string = "https://checkip.amazonaws.com/"

func GetMyIp() (*string, error) {
	resp, err := http.Get(AwsIpCheckerUrl)
	if err != nil {
		log.Printf("Cannot reach out to %s", AwsIpCheckerUrl)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Cannot read response from %s", AwsIpCheckerUrl)
		return nil, err
	}
	ip := string(body)
	return &ip, nil
}

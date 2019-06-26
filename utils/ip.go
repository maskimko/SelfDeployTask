package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

const AwsIpCheckerUrl string = "https://checkip.amazonaws.com/"

func GetMyIp() (*[]byte, error) {
	resp, err := http.Get(AwsIpCheckerUrl)
	if err != nil {
		log.Printf("Cannot reach out to %s\n", AwsIpCheckerUrl)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Cannot read response from %s\n", AwsIpCheckerUrl)
		return nil, err
	}
	ip := bytes.Trim(body, " \n")
	return &ip, nil
}

package config

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"gopkg.in/yaml.v3"
)

var (
	ServerURL string
	T         int
	N         int
	K         int
)

func Extract() (string, error) {
	if ServerURL != "" {
		return ServerURL, nil
	}
	yfile, err := ioutil.ReadFile("server-config.yaml")
	if err != nil {
		return "", err
	}
	configData := make(map[string]string)

	err = yaml.Unmarshal(yfile, &configData)
	if err != nil {
		return "", err
	}
	fmt.Println(configData)
	ServerURL = configData["server_url"]
	T, err = strconv.Atoi(configData["T"])
	if err != nil {
		return "", err
	}
	N, err = strconv.Atoi(configData["N"])
	if err != nil {
		return "", err
	}
	K, err = strconv.Atoi(configData["K"])
	if err != nil {
		return "", err
	}
	return configData["server_url"], nil
}

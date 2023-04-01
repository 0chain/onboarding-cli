package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var ServerURL string

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
	ServerURL = configData["server_url"]
	return configData["server_url"], nil
}

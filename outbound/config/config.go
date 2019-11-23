package config

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Config struct {
	AppId  string `yaml:"app_id"`
	AppKey string `yaml:"app_key"`
}

func GetConfig() (*Config, error) {
	configModel, err := doReadConfig("F:/webfiles/github/SchipholFlightAnalyser/outbound/config/config.yml") // TODO figure out why relative path does not work. Windows issue?
	if err != nil {
		return nil, err
	}

	return configModel, nil
}

func doReadConfig(filename string) (*Config, error) {
	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %s", filename, err)
	}

	configModel := Config{}
	err = yaml.Unmarshal([]byte(yamlData), &configModel)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling conf file %s: %s", filename, err)
	}

	return &configModel, nil
}

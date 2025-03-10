package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	MockPath = "%s/mock.config.yaml"
	RealPath = "%s/real.config.yaml"
)

type AccountInfo struct {
	ID          string `yaml:"ID"`
	Address     string `yaml:"Address"`
	ProductCode string `yaml:"ProductCode"`
}

type AppInfo struct {
	Key    string `yaml:"Key"`
	Secret string `yaml:"Secret"`
}

type AccessInfo struct {
	Token   string `yaml:"Token"`
	Expired string `yaml:"Expired"`
}

type Config struct {
	IsMock  bool        `yaml:"IsMock"`
	Url     string      `yaml:"Url"`
	Account AccountInfo `yaml:"Account"`
	App     AppInfo     `yaml:"App"`
	Access  AccessInfo  `yaml:"Access"`
}

func NewConfig(isMock bool, url string,
	id string, address string, code string, key string, secret string) *Config {
	return &Config{
		IsMock: isMock,
		Url:    url,
		Account: AccountInfo{
			Address:     address,
			ProductCode: code,
		},
		App: AppInfo{
			Key:    key,
			Secret: secret,
		},
	}
}

func getPath(isMock bool, path string) string {
	if isMock {
		return fmt.Sprintf(MockPath, path)
	} else {
		return fmt.Sprintf(RealPath, path)
	}
}

func LoadConfig(path string, isMock bool) (*Config, error) {
	var config Config

	data, _ := os.ReadFile(getPath(isMock, path))
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error while unmarshaling config:%v", err)
	} else {
		return &config, nil
	}
}

func (c Config) Save2File(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error while marshaling config:%v", err)
	}

	if err = os.WriteFile(getPath(c.IsMock, path), data, 0666); err != nil {
		return fmt.Errorf("error while writing to %s:%v", path, err)
	}

	return nil
}

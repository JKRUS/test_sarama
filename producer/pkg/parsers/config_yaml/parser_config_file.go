package config_yaml

import (
	"github.com/kylelemons/go-gypsy/yaml"
)

type config struct {
	address   string
	topicName string
}

func NewConfig() *config {
	return &config{
		address:   "",
		topicName: "",
	}
}

func (c *config) ReadFile(path string) ([]string, error) {
	config, err := yaml.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ipaddress, err := config.Get("ipAddress")
	if err != nil {
		return nil, err
	}
	port, err := config.Get("port")
	if err != nil {
		return nil, err
	}
	c.topicName, err = config.Get("topic")
	if err != nil {
		return nil, err
	}
	c.address = ipaddress + ":" + port
	res := []string{c.address, c.topicName}
	return res, nil
}

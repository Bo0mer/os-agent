package config

import (
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"
)

type OSAgentConfig struct {
	Id     string       `yaml:"id"`
	Host   string       `yaml:"host"`
	Port   int          `yaml:"port"`
	Master MasterConfig `yaml:"master"`
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MasterConfig struct {
	URL string `yaml:"url"`
}

func LoadConfig(configFile string) (OSAgentConfig, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return OSAgentConfig{}, err
	}

	config := &OSAgentConfig{}

	decoder := candiedyaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return OSAgentConfig{}, err
	}

	return *config, nil
}

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
	Host string     `yaml:"host"`
	Port int        `yaml:"port"`
	Auth AuthConfig `yaml:"auth"`
}

type MasterConfig struct {
	URL               string `yaml:"url"`
	HeartbeatInterval int    `yaml:"heartbeat_interval"`
}

type AuthConfig struct {
	Enabled  bool   `yaml:"enabled"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
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

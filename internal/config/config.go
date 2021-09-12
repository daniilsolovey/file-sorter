package config

import (
	"github.com/kovetskiy/ko"
	"gopkg.in/yaml.v1"
)

type Config struct {
	PathToInputDir  string `yaml:"path_to_input_dir" required:"true"`
	PathToResultDir string `yaml:"path_to_result_dir" required:"true"`
}

func Load(path string) (*Config, error) {
	config := &Config{}
	err := ko.Load(path, config, ko.RequireFile(false), yaml.Unmarshal)
	if err != nil {
		return nil, err
	}

	return config, nil
}

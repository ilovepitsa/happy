package config

import (
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	App struct {
		Name string `yaml:"name"`
	}

	Network struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Logger struct {
		JSONEnable bool   `yaml:"json_enable"`
		Level      string `yaml:"level"`
		Report     bool   `yaml:"report"`
	}

	Postgres struct {
		User     string `yaml:"host"`
		Password string `yaml:"password"`
		Db       string `yaml:"db"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	}

	Config struct {
		App App     `yaml:"app"`
		Net Network `yaml:"network"`
		Log Logger  `yaml:"logger"`
	}
)

func NewConfig(cnfPath string) (*Config, error) {
	config := &Config{}
	err := cleanenv.ReadConfig(path.Join("./", cnfPath), config)
	if err != nil {
		return nil, err
	}
	return config, nil

}

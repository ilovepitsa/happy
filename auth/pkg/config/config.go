package config

import (
	"path"
	"time"

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

	Redis struct {
		Host     string        `yaml:"host"`
		Port     string        `yaml:"port"`
		Password string        `yaml:"pass"`
		TTL      time.Duration `yaml:"ttl"`
	}

	Config struct {
		App App     `yaml:"app"`
		Net Network `yaml:"network"`
		Log Logger  `yaml:"logger"`
		R   Redis   `yaml:"redis"`
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

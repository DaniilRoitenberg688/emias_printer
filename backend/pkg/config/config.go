package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	Port int `yaml:"port" env:"PORT" env-default:8000`
	Delay int `yaml:"delay" env:"DELAY" env-default:5`
}


func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil

}

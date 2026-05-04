package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Conf *Config

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	JWT      JWTConfig      `yaml:"jwt"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type LoggerConfig struct {
	Level     string `yaml:"level"`
	Path      string `yaml:"path"`
	MaxSize   int    `yaml:"max_size"`
	MaxBackup int    `yaml:"max_backup"`
	MaxAge    int    `yaml:"max_age"`
}

func InitConfig() error {
	file, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return err
	}

	Conf = new(Config)
	return yaml.Unmarshal(file, Conf)
}

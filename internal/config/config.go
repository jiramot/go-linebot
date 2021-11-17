package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	LineConfig   LineConfig   `mapstructure:"line"`
	ServerConfig ServerConfig `mapstructure:"server"`
}

type LineConfig struct {
	ChannelSecret      string `mapstructure:"channel_secret"`
	ChannelAccessToken string `mapstructure:"channel_access_token"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.New("unable read config")
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.New("unable read config")
	}
	return &config, nil
}

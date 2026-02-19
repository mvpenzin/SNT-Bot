package main

import (
	"gopkg.in/ini.v1"
)

type Config struct {
	Telegram TelegramConfig
	Database DatabaseConfig
	Server   ServerConfig
	Weather  WeatherConfig
}

type TelegramConfig struct {
	Token string
	Debug bool
}

type DatabaseConfig struct {
	URL      string
	MaxConns int
}

type ServerConfig struct {
	Port int
}

type WeatherConfig struct {
	URL  string
	LAT  float64
	LON  float64
	PAST int
	DAYS int
	ZONE string
	WIND string
}

func LoadConfig(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	// Telegram
	config.Telegram.Token = cfg.Section("telegram").Key("token").String()
	config.Telegram.Debug = cfg.Section("telegram").Key("debug").MustBool(false)

	// Database
	config.Database.URL = cfg.Section("database").Key("url").String()
	config.Database.MaxConns = cfg.Section("database").Key("max_conns").MustInt(10)

	// Server
	config.Server.Port = cfg.Section("server").Key("port").MustInt(5000)

	// Weather
	config.Weather.URL = cfg.Section("weather").Key("url").String()
	config.Weather.LAT = cfg.Section("weather").Key("lat").MustFloat64(53.327935)
	config.Weather.LON = cfg.Section("weather").Key("lon").MustFloat64(84.102975)
	config.Weather.PAST = cfg.Section("weather").Key("past").MustInt(1)
	config.Weather.DAYS = cfg.Section("weather").Key("days").MustInt(3)
	config.Weather.ZONE = cfg.Section("weather").Key("zone").String()
	config.Weather.WIND = cfg.Section("weather").Key("wind").String()

	return config, nil
}

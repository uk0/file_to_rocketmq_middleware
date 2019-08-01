package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

type tomlConfig struct {
	FtpSender FtpSender
}
type FtpSender struct {
	WatchDir           string `toml:"watchDir"`
	RocketMQNameserver string `toml:"rocketMQNameserver"`
	SendTopic          string `toml:"sendTopic"`
	SendGroup          string `toml:"sendGroup"`
}

func GetConfig() FtpSender {
	var favorites tomlConfig
	if _, err := toml.DecodeFile("config.toml", &favorites); err != nil {
		log.Fatal(err)
	}
	return favorites.FtpSender;
}

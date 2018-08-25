package server

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type Config struct {
	DiscordApi       string `required:"true" split_words:"true"`
	ChallongeApi     string `split_words:"true"`
	Subdomain        string `desc:"Challonge subdomain"`
	ConsumerKey      string `desc:"Twitter consumer key" split_words:"true"`
	ConsumerSecret   string `desc:"Twitter consumer secret" split_words:"true"`
	AccessToken      string `desc:"Twitter access token" split_words:"true"`
	AccessSecret     string `desc:"Twitter access secret" split_words:"true"`
	PostChannel      string `desc:"channel id to post" split_words:"true"`
	AdminChannel     string `desc:"channel id to post errors" split_words:"true"`
	Database         string `desc:"backend postgres database"`
	DatabaseHost     string `desc:"backend postgres host" split_words:"true"`
	DatabaseUser     string `desc:"backend user" split_words:"true" default:"postgres"`
	DatabasePassword string `desc:"backend password" split_words:"true"`
}

var config Config

func GetConfig() Config {
	return config
}

func SetConfig() {
	envconfig.Usage("discord_bot", &config)
	if err := envconfig.Process("discord_bot", &config); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

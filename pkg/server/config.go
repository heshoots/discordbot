package server

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"

	"github.com/heshoots/discordbot/pkg/config"
	"github.com/heshoots/discordbot/pkg/models"
)

type Config struct {
	DiscordApi     string `required:"true" split_words:"true"`
	ChallongeApi   string `split_words:"true"`
	Subdomain      string `desc:"Challonge subdomain"`
	ConsumerKey    string `desc:"Twitter consumer key" split_words:"true"`
	ConsumerSecret string `desc:"Twitter consumer secret" split_words:"true"`
	AccessToken    string `desc:"Twitter access token" split_words:"true"`
	AccessSecret   string `desc:"Twitter access secret" split_words:"true"`
	PostChannel    string `desc:"channel id to post" split_words:"true"`
	AdminChannel   string `desc:"channel id to post errors" split_words:"true"`
}

var appconfig Config

func GetConfig() Config {
	return appconfig
}

func AddRole(role *models.Role) error {
	return config.AddRole(role)
}

func SetConfig() {
	config.ReadConfigFile("config/config.json")
	envconfig.Usage("discord_bot", &appconfig)
	if err := envconfig.Process("discord_bot", &appconfig); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"github.com/heshoots/discordbot/models"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var config struct {
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

var compiled string

func main() {
	envconfig.Usage("discord_bot", &config)
	if err := envconfig.Process("discord_bot", &config); err != nil {
		log.Println(os.Stderr, err)
		os.Exit(1)
	}
	models.DB(config.DatabaseHost, config.Database, config.DatabaseUser, config.DatabasePassword)
	discord, err := NewRouter(config.DiscordApi)
	if err != nil {
		log.Panic(err)
		return
	}
	err = discord.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}
	discord.ChannelMessageSend(config.AdminChannel, "Redeployed, compiled: "+compiled)
	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/kelseyhightower/envconfig"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var config struct {
        DiscordApi     string `required:"true" split_words:"true"`
        ConsumerKey    string `desc:"Twitter consumer key" split_words:"true"`
        ConsumerSecret string `desc:"Twitter consumer secret" split_words:"true"`
        AccessToken    string `desc:"Twitter access token" split_words:"true"`
        AccessSecret   string `desc:"Twitter access secret" split_words:"true"`
        PostChannel    string `desc:"channel id to post" split_words:"true"`
}

func main() {
	if err := envconfig.Process("discord_bot", &config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		envconfig.Usage("discord_bot", &config)
		os.Exit(1)
	}
	discord, err := discordgo.New("Bot " + config.DiscordApi)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}
	discord.AddHandler(messageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func tweet(message string) {
	fmt.Println("tweeting")
	twitterconfig := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	token := oauth1.NewToken(config.AccessToken, config.AccessSecret)
	httpClient := twitterconfig.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	_, _, err := client.Statuses.Update(message, nil)
	if err != nil {
		fmt.Println("could not tweet,", err)
	}
}

func runCommand(s *discordgo.Session, command string, message string) {
	fmt.Println(command, "..", message)
	if command == "!twitter " {
		fmt.Println("match")
	}
	switch command {
	case "!twitter":
		tweet(message)
	case "!discord":
		s.ChannelMessageSend(config.PostChannel, message)
	default:
		fmt.Println("didn't match")
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	const param string = "!"
	if m.Content[0:1] == param {
		split := strings.SplitAfterN(m.Content, " ", 2)
		command := strings.Trim(split[0], " ")
		message := strings.Trim(split[1], " ")
		runCommand(s, command, message)
	}
}

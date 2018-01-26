package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var config struct {
	DiscordApi     string `required:"true" split_words:"true"`
	ChallongeApi   string `split_words:"true"`
	ConsumerKey    string `desc:"Twitter consumer key" split_words:"true"`
	ConsumerSecret string `desc:"Twitter consumer secret" split_words:"true"`
	AccessToken    string `desc:"Twitter access token" split_words:"true"`
	AccessSecret   string `desc:"Twitter access secret" split_words:"true"`
	PostChannel    string `desc:"channel id to post" split_words:"true"`
}

func main() {
	envconfig.Usage("discord_bot", &config)
	if err := envconfig.Process("discord_bot", &config); err != nil {
		fmt.Fprintln(os.Stderr, err)
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

func runCommand(s *discordgo.Session, command string, message string) error {
	fmt.Println(command, "..", message)
	if command == "!twitter " {
		fmt.Println("match")
	}
	switch command {
	case "!announce":
		tweet(message)
		s.ChannelMessageSend(config.PostChannel, message)
	case "!twitter":
		tweet(message)
	case "!discord":
		s.ChannelMessageSend(config.PostChannel, message)
	case "!challonge":
		split := strings.SplitAfterN(message, " ", 2)
		name := strings.Trim(split[0], " ")
		game := strings.Trim(split[1], " ")
		url, err := createTournament(name, game)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(config.PostChannel, url)
	default:
		return errors.New("command not recognised")
	}
	return nil
}

func createTournament(name string, game string) (string, error) {
	client := &http.Client{}
	tournamentvalues := map[string]string{"name": name, "url": name, "subdomain": "smbf", "game_name": game}
	values := map[string]map[string]string{"tournament": tournamentvalues}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.challonge.com/v1/tournaments.json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("api_key", config.ChallongeApi)
	req.URL.RawQuery = q.Encode()
	_, err = client.Do(req)
	if err != nil {
		return "", err
	}
	return "http://smbf.challonge.com/" + name, nil
}

func isAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		fmt.Println("could not access permissions, ", err)
	}
	return permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	const param string = "!"
	if isAdmin(s, m) && m.Content[0:1] == param {
		split := strings.SplitAfterN(m.Content, " ", 2)
		command := strings.Trim(split[0], " ")
		var message string = ""
		if len(split) > 1 {
			message = strings.Trim(split[1], " ")
		}
		err := runCommand(s, command, message)
		if err != nil {
			fmt.Println("Failed to run command", err)
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		}
	}
}

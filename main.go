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
	AdminChannel   string `desc:"channel id to post errors" split_words:"true"`
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

func adminCommand(s *discordgo.Session, command string, message string) error {
	fmt.Println(command, message)
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
		if len(split) != 2 {
			return errors.New("not enough input, command: !challonge url game_name")
		}
		name := strings.Trim(split[0], " ")
		game := strings.Trim(split[1], " ")
		url, err := createTournament(name, game)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(config.PostChannel, url)
	}
	return nil
}

func userCommand(s *discordgo.Session, command string, message string, channel string) {
	switch command {
	case "!hi":
		s.ChannelMessageSend(channel, "https://78.media.tumblr.com/c52387b2f0599b6aad20defb9b3ad6b9/tumblr_ngwarrlkfG1qcm0i5o2_500.gif")
	}
}

func createTournament(name string, game string) (string, error) {
	client := &http.Client{}
	tournamentvalues := map[string]string{"name": name, "url": name, "subdomain": "smbf", "game_name": game, "tournament_type": "double elimination"}
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
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return "", errors.New(resp.Status + "challonge create failed " + buf.String())
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
	if m.Content[0:1] == param {
		split := strings.SplitAfterN(m.Content, " ", 2)
		command := strings.Trim(split[0], " ")
		var message string = ""
		if len(split) > 1 {
			message = strings.Trim(split[1], " ")
		}
		if isAdmin(s, m) {
			err := adminCommand(s, command, message)
			if err != nil {
				fmt.Println("Failed to run command", err)
				s.ChannelMessageSend(config.AdminChannel, "Error: "+err.Error())
			}
		}
		userCommand(s, command, message, m.ChannelID)
	}
}

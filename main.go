package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/heshoots/discordbot/models"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

func tweet(message string) string {
	log.Println("tweeting")
	twitterconfig := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	token := oauth1.NewToken(config.AccessToken, config.AccessSecret)
	httpClient := twitterconfig.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	tweet, _, err := client.Statuses.Update(message, nil)
	if err != nil {
		log.Println("could not tweet,", err)
	}
	return "http://twitter.com/sodiumshowdown/status/" + fmt.Sprint(tweet.ID)
}

func createTournament(name string, game string) (string, error) {
	client := &http.Client{}
	tournamentvalues := map[string]string{"name": name, "url": name, "subdomain": config.Subdomain, "game_name": game, "tournament_type": "double elimination"}
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
	return "http://" + config.Subdomain + ".challonge.com/" + name, nil
}

func isAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		log.Println("could not access permissions, ", err)
	}
	return permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator
}

func hasPrefix(prefix string, m *discordgo.MessageCreate) bool {
	if len(m.Content) >= len(prefix) {
		return m.Content[0:len(prefix)] == prefix
	}
	return false
}

func getCommand(m *discordgo.MessageCreate) string {
	split := strings.SplitAfterN(m.Content, " ", 2)
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

func getGuild(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.Guild, error) {
	// Attempt to get the channel from the state.
	// If there is an error, fall back to the restapi
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		channel, err = s.Channel(m.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	// Attempt to get the guild from the state,
	// If there is an error, fall back to the restapi.
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		guild, err = s.Guild(channel.GuildID)
		if err != nil {
			return nil, err
		}
	}
	return guild, nil
}

func getRoles(s *discordgo.Session, m *discordgo.MessageCreate) ([]*discordgo.Role, error) {
	guild, err := getGuild(s, m)
	if err != nil {
		log.Println("couldn't get guildID, ", err)
		return nil, err
	}
	roles, err := s.GuildRoles(guild.ID)
	if err != nil {
		log.Println("couldn't get roles, ", err)
		return nil, err
	}
	return roles, err
}

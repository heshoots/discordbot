package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/kelseyhightower/envconfig"
	"github.com/heshoots/discordbot/models"
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
}

var compiled string

func main() {
	envconfig.Usage("discord_bot", &config)
	if err := envconfig.Process("discord_bot", &config); err != nil {
		log.Println(os.Stderr, err)
		os.Exit(1)
	}
	models.DB(config.DatabaseHost, config.Database)
	discord, err := discordgo.New("Bot " + config.DiscordApi)
	if err != nil {
		log.Println("error creating Discord session,", err)
	}
	discord.AddHandler(prefixHandler("!discord", discordHandler))
	discord.AddHandler(prefixHandler("!announce", discordHandler))

	discord.AddHandler(prefixHandler("!challonge", challongeHandler))

	discord.AddHandler(prefixHandler("!twitter", twitterHandler))
	discord.AddHandler(prefixHandler("!tweet", twitterHandler))
	discord.AddHandler(prefixHandler("!announce", twitterHandler))

	// Discord Role Handlers
	discord.AddHandler(prefixHandler("!makerole", makeRoleHandler))
	discord.AddHandler(prefixHandler("!removerole", removeRoleHandler))
	discord.AddHandler(prefixHandler("!showroles", showRolesHandler))
	discord.AddHandler(prefixHandler("!giverole", giveRoleHandler))
	discord.AddHandler(prefixHandler("!takerole", takeRoleHandler))

	// Fun handler
	discord.AddHandler(prefixHandler("!hi", hiHandler))

	err = discord.Open()
	if err != nil {
		log.Println("error opening connection,", err)
	}
	discord.ChannelMessageSend(config.AdminChannel, "Redeployed, compiled: " + compiled)
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
	return "http://" +  config.Subdomain + ".challonge.com/" + name, nil
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
	if (len(split) > 1) {
		return split[1]
	}
	return ""
}

func prefixHandler(prefix string, handler func (*discordgo.Session, *discordgo.MessageCreate)) func (s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if hasPrefix(prefix, m) {
			handler(s, m)
		}
	}
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

func getRoles(s *discordgo.Session, m *discordgo.MessageCreate) ([]*discordgo.Role, error)  {
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

func removeRoleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if isAdmin(s, m) {
		command := getCommand(m)
		err := models.DeleteRole(command)
		if err != nil {
			log.Println("couldn't delete role, ", err)
			s.ChannelMessageSend(config.AdminChannel, "couldn't delete role")
			return
		}
		s.ChannelMessageSend(config.AdminChannel, "role deleted")
	}
}

func makeRoleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if isAdmin(s, m) {
		roles, err := getRoles(s, m)
		if err != nil {
			s.ChannelMessageSend(config.AdminChannel, "couldn't create role")
			log.Println("couldn't get roles, ", err)
			return
		}
		command := getCommand(m)
		for _, role := range roles {
			if command == role.Name {
				role := models.Role{Name: role.Name, RoleID: role.ID}
				err := models.CreateRole(&role)
				if err != nil {
					s.ChannelMessageSend(config.AdminChannel, "couldn't create role")
					log.Println("couldn't create role, ", err)
					return

				} else {
					s.ChannelMessageSend(config.AdminChannel, "Role added")
					return
				}
			}
		}
	}
}

func showRolesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	rolesHelp :=  `
To get a role use !giverole Role
To remove a role use !takerole Role

Available Roles
-----------
`
	roles, err := models.GetRoles()
	if err != nil {
		log.Println("couldn't show roles, ", err)
		s.ChannelMessageSend(config.PostChannel, "couldn't show roles")
		return
	}
	var output string
	for _, role := range roles {
		output = output + "\n" + role.Name
	}
	s.ChannelMessageSend(m.ChannelID, "``` " + rolesHelp + output + " ```" )
}

func giveRoleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := getCommand(m)
	guild, _ := getGuild(s, m)
	if command == "Thems Fighting Herds" {
		err := s.GuildBanCreateWithReason(guild.ID, m.Author.ID, "Furry", 0)
		if err != nil {
			log.Println("Something went wrong", err)
			s.ChannelMessageSend(m.ChannelID, "Something went wrong")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "BANNED: get owned")
		return
	}
	role, err := models.GetRole(command)
	if err != nil {
		log.Println("Role unavailable", err)
		s.ChannelMessageSend(m.ChannelID, "couldn't add role")
		return
	}
	err = s.GuildMemberRoleAdd(guild.ID, m.Author.ID, role.RoleID)
	if err != nil {
		log.Println("couldn't add role", err)
		s.ChannelMessageSend(m.ChannelID, "couldn't add role")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Role added")
}

func takeRoleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := getCommand(m)
	role, err := models.GetRole(command)
	if err != nil {
		log.Println("Role unavailable", err)
		s.ChannelMessageSend(m.ChannelID, "couldn't remove role")
		return
	}
	guild, _ := getGuild(s, m)
	err = s.GuildMemberRoleRemove(guild.ID, m.Author.ID, role.RoleID)
	if err != nil {
		log.Println("couldn't remove role", err)
		s.ChannelMessageSend(m.ChannelID, "couldn't remove role")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Role removed")
}

func discordHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if isAdmin(s, m) {
		s.ChannelMessageSend(config.PostChannel, getCommand(m))
	}
}

func challongeHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if isAdmin(s, m) {
		command := getCommand(m)
		split := strings.SplitAfterN(command , " ", 2)
		if len(split) != 2 {
			s.ChannelMessageSend(config.AdminChannel, "not enough input, command: !challonge url game_name")
			return
		}
		name := strings.Trim(split[0], " ")
		game := strings.Trim(split[1], " ")
		url, err := createTournament(name, game)
		if err != nil {
			s.ChannelMessageSend(config.AdminChannel, "couldn't create tournament: " + err.Error())
			return
		}
		s.ChannelMessageSend(config.PostChannel, url)
	}

}

func twitterHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if isAdmin(s, m) {
		url := tweet(getCommand(m))
		s.ChannelMessageSend(config.AdminChannel, url)
	}
}

func hiHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "https://78.media.tumblr.com/c52387b2f0599b6aad20defb9b3ad6b9/tumblr_ngwarrlkfG1qcm0i5o2_500.gif")
}

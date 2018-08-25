package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/heshoots/discordbot/discordhelpers"
	"github.com/heshoots/discordbot/events"
	"github.com/heshoots/discordbot/models"
	"github.com/heshoots/discordbot/twitter"
	"log"
	"time"
)

func isAdminHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate)) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if discordhelpers.IsAdmin(s, m) {
			handler(s, m)
		} else {
			s.ChannelMessageSend(m.ChannelID, "You dont have permissions to do that")
		}
		return
	}
}

func prefixHandler(prefix string, handler func(*discordgo.Session, *discordgo.MessageCreate)) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if discordhelpers.HasPrefix(prefix, m) {
			handler(s, m)
		}
	}
}

func removeRoleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if discordhelpers.IsAdmin(s, m) {
		command := discordhelpers.GetCommand(m)
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
	if discordhelpers.IsAdmin(s, m) {
		roles, err := discordhelpers.GetRoles(s, m)
		if err != nil {
			s.ChannelMessageSend(config.AdminChannel, "couldn't create role")
			log.Println("couldn't get roles, ", err)
			return
		}
		command := discordhelpers.GetCommand(m)
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

func inviteHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Invite Link: <http://discord.superminerbattle.farm>")
}

func showRolesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	rolesHelp := `
To get a role use !iam Role
To remove a role use !iamn Role

Roles ending in "Fighters" can be @ mentioned

Available Roles
-----------
`
	roles, err := models.GetRoles()
	if err != nil {
		log.Println("couldn't show roles, ", err)
		s.ChannelMessageSend(m.ChannelID, "couldn't show roles")
		return
	}
	var output string
	for _, role := range roles {
		output = output + "\n !iam " + role.Name
	}
	s.ChannelMessageSend(m.ChannelID, "``` "+rolesHelp+output+" ```")
}

func iamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := discordhelpers.GetCommand(m)
	guild, _ := discordhelpers.GetGuild(s, m)
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

func iamnHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := discordhelpers.GetCommand(m)
	role, err := models.GetRole(command)
	if err != nil {
		log.Println("Role unavailable", err)
		s.ChannelMessageSend(m.ChannelID, "couldn't remove role")
		return
	}
	guild, _ := discordhelpers.GetGuild(s, m)
	err = s.GuildMemberRoleRemove(guild.ID, m.Author.ID, role.RoleID)
	if err != nil {
		log.Println("couldn't remove role", err)
		s.ChannelMessageSend(m.ChannelID, "couldn't remove role")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Role removed")
}

func discordHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if discordhelpers.IsAdmin(s, m) {
		if discordhelpers.HasPrefix("!announce", m) {
			s.ChannelMessageSend(config.PostChannel, "@everyone "+discordhelpers.GetCommand(m))
		} else {
			s.ChannelMessageSend(config.PostChannel, discordhelpers.GetCommand(m))
		}
	}
}

func twitterHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if discordhelpers.IsAdmin(s, m) {
		auth := twitter.TwitterAuth{
			config.ConsumerKey,
			config.ConsumerSecret,
			config.AccessToken,
			config.AccessSecret,
		}
		url, err := twitter.Tweet(auth, discordhelpers.GetCommand(m))
		if err != nil {
			s.ChannelMessageSend(config.AdminChannel, err.Error())
		}
		s.ChannelMessageSend(config.AdminChannel, url)
	}
}

func helpHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	var help string
	if m.ChannelID == config.AdminChannel {
		for _, route := range GetRoutes() {
			help += route.Prefix[0] + " : " + route.HelpText + "\n"
		}
	} else {
		for _, route := range GetRoutes() {
			if !route.Admin {
				help += route.Prefix[0] + " : " + route.HelpText + "\n"
			}
		}
	}
	s.ChannelMessageSend(m.ChannelID, help)
}

func Logger(route Route) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		start := time.Now()
		route.Handler(s, m)
		log.Printf(
			"%s\t%t\t%s\t%s\t%s",
			route.Name,
			route.Admin,
			m.Author.Username,
			m.Content,
			time.Since(start),
		)
	}
}

func hiHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "https://78.media.tumblr.com/c52387b2f0599b6aad20defb9b3ad6b9/tumblr_ngwarrlkfG1qcm0i5o2_500.gif")
}

func lanesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	eventlist, err := events.GetLanesEvents()
	if err != nil {
		log.Fatal(err)
		return
	}
	out := ""
	for _, event := range eventlist {
		if event.Date.Weekday() == time.Sunday {
			out = out + event.Title + "\n" + event.Date.Format("Mon Jan 2") + "\n" + event.Description + "\n\n"
		}
	}
	s.ChannelMessageSend(m.ChannelID, out)
}

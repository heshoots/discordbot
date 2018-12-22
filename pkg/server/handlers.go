package server

import (
	"github.com/bwmarrin/discordgo"
	"github.com/heshoots/discordbot/pkg/discord"
	"github.com/heshoots/discordbot/pkg/events"
	"github.com/heshoots/discordbot/pkg/models"
	"github.com/heshoots/discordbot/pkg/twitter"
	"log"
	"time"
)

func isAdminHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate)) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if discord.IsAdmin(s, m) {
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
		if discord.HasPrefix(prefix, m) {
			handler(s, m)
		}
	}
}

func statusHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := discord.GetCommand(m)
	s.UpdateStatus(0, command)
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
	roles, err := models.YamlRoles()
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
	command := discord.GetCommand(m)
	guild, _ := discord.GetGuild(s, m)
	role, err := models.YamlRole(command)
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
	command := discord.GetCommand(m)
	role, err := models.YamlRole(command)
	if err != nil {
		log.Println("Role unavailable", err)
		s.ChannelMessageSend(m.ChannelID, "couldn't remove role")
		return
	}
	guild, _ := discord.GetGuild(s, m)
	err = s.GuildMemberRoleRemove(guild.ID, m.Author.ID, role.RoleID)
	if err != nil {
		log.Println("couldn't remove role", err)
		s.ChannelMessageSend(m.ChannelID, "couldn't remove role")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Role removed")
}

func discordHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if discord.IsAdmin(s, m) {
		if discord.HasPrefix("!announce", m) {
			s.ChannelMessageSend(GetConfig().PostChannel, "@everyone "+discord.GetCommand(m))
		} else {
			s.ChannelMessageSend(GetConfig().PostChannel, discord.GetCommand(m))
		}
	}
}

func twitterHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if discord.IsAdmin(s, m) {
		auth := twitter.TwitterAuth{
			GetConfig().ConsumerKey,
			GetConfig().ConsumerSecret,
			GetConfig().AccessToken,
			GetConfig().AccessSecret,
		}
		url, err := twitter.Tweet(auth, discord.GetCommand(m))
		if err != nil {
			s.ChannelMessageSend(GetConfig().AdminChannel, err.Error())
		}
		s.ChannelMessageSend(GetConfig().AdminChannel, url)
	}
}

func helpHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	var help string
	if m.ChannelID == GetConfig().AdminChannel {
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

func roleNameAddedHandler(s *discordgo.Session, role *discordgo.Role) {
	time.Sleep(10 * time.Second)
	log.Println("In function")
	s.ChannelMessageSend(GetConfig().AdminChannel, "Looks like you just created a role "+role.Mention()+" add it to available roles? Use `!addrole "+role.ID+"`")
}

func roleAddedHandler(s *discordgo.Session, m *discordgo.GuildRoleCreate) {
	guildRole := m.GuildRole
	go roleNameAddedHandler(s, guildRole.Role)
}

func AddRoleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(GetConfig().PostChannel, discord.GetCommand(m))
	roles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "couldn't find role")
	}
	for _, role := range roles {
		if role.ID == discord.GetCommand(m) {
			err = AddRole(&models.Role{Name: role.Name, RoleID: role.ID})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "couldn't add role to database")
			}
			s.ChannelMessageSend(m.ChannelID, "role added")
		}
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

package discordhelpers

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func IsAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		log.Println("could not access permissions, ", err)
	}
	return permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator
}

func HasPrefix(prefix string, m *discordgo.MessageCreate) bool {
	if len(m.Content) >= len(prefix) {
		return m.Content[0:len(prefix)] == prefix
	}
	return false
}

func GetCommand(m *discordgo.MessageCreate) string {
	split := strings.SplitAfterN(m.Content, " ", 2)
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

func GetGuild(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.Guild, error) {
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

func GetRoles(s *discordgo.Session, m *discordgo.MessageCreate) ([]*discordgo.Role, error) {
	guild, err := GetGuild(s, m)
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

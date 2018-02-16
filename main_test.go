package main

import (
	"github.com/bwmarrin/discordgo"
	"testing"
)

type MockMessage struct {
	content        string
	authorID       string
	authorUsername string
	channelID      string
}

func (m *MockMessage) GetContent() string {
	return m.content
}

func (m *MockMessage) GetAuthorID() string {
	return m.authorID
}

func (m *MockMessage) GetAuthorUsername() string {
	return m.authorUsername
}

func (m *MockMessage) GetChannelID() string {
	return m.channelID
}

func (m *MockMessage) addUserId(userID string) *MockMessage {
	m.authorID = userID
	return m
}

func (m *MockMessage) SetContent(content string) *MockMessage {
	m.content = content
	return m
}

func (m *MockMessage) SetChannelID(channelID string) *MockMessage {
	m.channelID = channelID
	return m
}

type MockSession struct {
	UserID    string
	Buffer    *OutMessage
	Admin     bool
	BotUserID string
}

type OutMessage struct {
	Channel string
	Message string
}

func (s *MockSession) ChannelMessageSend(channel string, message string) (*discordgo.Message, error) {
	s.Buffer = &OutMessage{channel, message}
	return &discordgo.Message{Content: message, ChannelID: message}, nil
}

func (s *MockSession) BotID() string {
	return s.BotUserID
}

func (s *MockSession) GuildMemberRoleAdd(string, string, string) error {
	return nil
}

func (s *MockSession) GuildMemberRoleRemove(string, string, string) error {
	return nil
}

func (s *MockSession) GuildRoles(string) ([]*discordgo.Role, error) {
	return nil, nil
}

func (s *MockSession) UserChannelPermissions(string, string) (int, error) {
	if s.Admin {
		return discordgo.PermissionAdministrator, nil
	}
	return 0, nil
}

func (s *MockSession) IsAdmin(string, string) bool {
	return s.Admin
}

func (s *MockSession) Guild(channel string) (*discordgo.Guild, error) {
	return &discordgo.Guild{}, nil
}

func (s *MockSession) SetAdmin(admin bool) *MockSession {
	s.Admin = admin
	return s
}

func (s *MockSession) SetBotID(botID string) *MockSession {
	s.BotUserID = botID
	return s
}

func (s *MockSession) setUserId(userID string) *MockSession {
	s.UserID = userID
	return s
}

func TestHasPrefix(t *testing.T) {
	msg := &MockMessage{}
	msg.addUserId("10").SetContent("!hello world")
	res := hasPrefix("!hello", msg)
	if res != true {
		t.Error("Message contains prefix")
	}
	msg.addUserId("10").SetContent("no !hello world")
	res = hasPrefix("!hello", msg)
	if res != false {
		t.Error("Prefix is not at beginning of message")
	}
	msg.addUserId("10").SetContent("!test world")
	res = hasPrefix("!hello", msg)
	if res != false {
		t.Error("Prefix does not occur")
	}
	msg.addUserId("10").SetContent("!hell")
	res = hasPrefix("!hello", msg)
	if res != false {
		t.Error("Message is shorter than prefix")
	}
}

func TestGetCommand(t *testing.T) {
	msg := &MockMessage{}
	msg.addUserId("10").SetContent("!command this is a command")
	command := getCommand(msg)
	if command != "this is a command" {
		t.Error("Returned the incorrect command")
	}
	msg.addUserId("10").SetContent("!nocommand")
	command = getCommand(msg)
	if command != "" {
		t.Error("didn't return empty command")
	}
	msg.addUserId("10").SetContent("!nocommand ")
	command = getCommand(msg)
	if command != "" {
		t.Error("doesn't handle extra spaces in command")
	}
}

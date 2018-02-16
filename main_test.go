package main

import (
	"github.com/bwmarrin/discordgo"
	"testing"
)

type MockMessage struct {
	*discordgo.Message
}

func (m *MockMessage) addUserId(userID string) *MockMessage {
	author := new(discordgo.User)
	author.ID = userID
	m.Author = author
	return m
}

func (m *MockMessage) setContent(content string) *MockMessage {
	m.Content = content
	return m
}

func (m *MockMessage) getMessageCreate() *discordgo.MessageCreate {
	return &discordgo.MessageCreate{m.Message}
}

type MockSession struct {
	*discordgo.Session
	buffer []*OutMessage
}

type OutMessage struct {
	Channel string
	Message string
}

func (s *MockSession) setUserId(userId string) *MockSession {
	s.State = new(discordgo.State)
	s.State.User = new(discordgo.User)
	s.State.User.ID = userId
	return s
}

func (s *MockSession) ChannelMessageSend(channel string, message string) {
	s.buffer = append(s.buffer, &OutMessage{channel, message})
}

type MockState struct {
	*discordgo.State
}

func (*MockState) isAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return true
}

func TestHasPrefix(t *testing.T) {
	msg := &MockMessage{&discordgo.Message{}}
	msg.addUserId("10").setContent("!hello world")
	res := hasPrefix("!hello", msg.getMessageCreate())
	if res != true {
		t.Error("Message contains prefix")
	}
	msg.addUserId("10").setContent("no !hello world")
	res = hasPrefix("!hello", msg.getMessageCreate())
	if res != false {
		t.Error("Prefix is not at beginning of message")
	}
	msg.addUserId("10").setContent("!test world")
	res = hasPrefix("!hello", msg.getMessageCreate())
	if res != false {
		t.Error("Prefix does not occur")
	}
	msg.addUserId("10").setContent("!hell")
	res = hasPrefix("!hello", msg.getMessageCreate())
	if res != false {
		t.Error("Message is shorter than prefix")
	}
}

func TestGetCommand(t *testing.T) {
	msg := &MockMessage{&discordgo.Message{}}
	msg.addUserId("10").setContent("!command this is a command")
	command := getCommand(msg.getMessageCreate())
	if command != "this is a command" {
		t.Error("Returned the incorrect command")
	}
	msg.addUserId("10").setContent("!nocommand")
	command = getCommand(msg.getMessageCreate())
	if command != "" {
		t.Error("didn't return empty command")
	}
	msg.addUserId("10").setContent("!nocommand ")
	command = getCommand(msg.getMessageCreate())
	if command != "" {
		t.Error("doesn't handle extra spaces in command")
	}
}

package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/heshoots/discordbot/discordhelpers"
	"testing"
)

func GetMockSession(userid string) *discordgo.Session {
	var session *discordgo.Session
	session = new(discordgo.Session)
	session.State = new(discordgo.State)
	session.State.User = new(discordgo.User)
	session.State.User.ID = userid
	return session
}

func GetMockMessageCreate(content string, authorID string) *discordgo.MessageCreate {
	author := new(discordgo.User)
	author.ID = authorID
	var message *discordgo.Message
	message = new(discordgo.Message)
	message.Author = author
	message.Content = content
	var messageEvent *discordgo.MessageCreate
	messageEvent = new(discordgo.MessageCreate)
	messageEvent.Message = message
	return messageEvent
}

func TestHasPrefix(t *testing.T) {
	msg := GetMockMessageCreate("!hello world", "10")
	res := discordhelpers.HasPrefix("!hello", msg)
	if res != true {
		t.Error("Message contains prefix")
	}
	msg = GetMockMessageCreate("no !hello world", "10")
	res = discordhelpers.HasPrefix("!hello", msg)
	if res != false {
		t.Error("Prefix is not at beginning of message")
	}
	msg = GetMockMessageCreate("!test world", "10")
	res = discordhelpers.HasPrefix("!hello", msg)
	if res != false {
		t.Error("Prefix does not occur")
	}
	msg = GetMockMessageCreate("!hell", "10")
	res = discordhelpers.HasPrefix("!hello", msg)
	if res != false {
		t.Error("Message is shorter than prefix")
	}
}

func TestPrefixHandler(t *testing.T) {
	session := GetMockSession("10")
	messageCreate := GetMockMessageCreate("!nocall", "10")
	prefixed := prefixHandler("!nocall", func(s *discordgo.Session, m *discordgo.MessageCreate) {
		t.Error("Should not be called")
	})
	prefixed(session, messageCreate)
	session = GetMockSession("11")
	messageCreate = GetMockMessageCreate("!call", "10")
	called := false
	prefixed = prefixHandler("!call", func(s *discordgo.Session, m *discordgo.MessageCreate) {
		called = true
	})
	prefixed(session, messageCreate)
	if called == false {
		t.Error("Handler should have been called")
	}
}

func TestGetComman(t *testing.T) {
	messageCreate := GetMockMessageCreate("!command this is a command", "10")
	command := discordhelpers.GetCommand(messageCreate)
	if command != "this is a command" {
		t.Error("Returned the incorrect command")
	}
	messageCreate = GetMockMessageCreate("!nocommand", "10")
	command = discordhelpers.GetCommand(messageCreate)
	if command != "" {
		t.Error("didn't return empty command")
	}
	messageCreate = GetMockMessageCreate("!nocommand ", "10")
	command = discordhelpers.GetCommand(messageCreate)
	if command != "" {
		t.Error("doesn't handle extra spaces in command")
	}
}

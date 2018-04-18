package main

import (
	"github.com/bwmarrin/discordgo"
)

type HandlerFunc func(s *discordgo.Session, m *discordgo.MessageCreate) bool

type Route struct {
	Name    string
	Prefix  []string
	Handler func(s *discordgo.Session, m *discordgo.MessageCreate)
	Admin   bool
}

type Routes []Route

var routes = Routes{
	Route{
		"Announce",
		[]string{"!discord", "!announce"},
		discordHandler,
		true,
	},
	Route{
		"Challonge",
		[]string{"!challonge"},
		challongeHandler,
		true,
	},
	Route{
		"Twitter",
		[]string{"!twitter", "!tweet", "!announce"},
		twitterHandler,
		true,
	},
	Route{
		"Make Role",
		[]string{"!makerole"},
		makeRoleHandler,
		true,
	},
	Route{
		"Remove Role",
		[]string{"!removerole"},
		removeRoleHandler,
		true,
	},
	Route{
		"Show Roles",
		[]string{"!showroles"},
		showRolesHandler,
		true,
	},
	Route{
		"Give Role",
		[]string{"!iam "},
		iamHandler,
		false,
	},
	Route{
		"Take Role",
		[]string{"!iamn"},
		iamnHandler,
		false,
	},
	Route{
		"Logging",
		[]string{"!"},
		loggingHandler,
		false,
	},
	Route{
		"Role Call",
		[]string{""},
		RoleCallHandler,
		false,
	},
	Route{
		"Help Handler",
		[]string{"!help"},
		RoleCallHandler,
		false,
	},
	Route{
		"Chloe hi",
		[]string{"!hi"},
		hiHandler,
		false,
	},
}

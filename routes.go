package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/heshoots/discordbot/challonge"
)

type HandlerFunc func(s *discordgo.Session, m *discordgo.MessageCreate)

type Route struct {
	Name    string
	Prefix  []string
	Handler HandlerFunc
	Admin   bool
}

type Routes []Route

func GetRoutes() Routes {
	return Routes{
		Route{
			"Announce",
			[]string{"!discord", "!announce"},
			discordHandler,
			true,
		},
		Route{
			"Challonge",
			[]string{"!challonge"},
			challonge.ChallongeHandler(config.ChallongeApi, config.Subdomain, []string{config.AdminChannel, config.PostChannel}, []string{config.AdminChannel}),
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
}

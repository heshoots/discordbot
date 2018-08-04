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
	Logged  bool
}

type Routes []Route

func GetRoutes() Routes {
	return Routes{
		Route{
			"Announce",
			[]string{"!discord", "!announce"},
			discordHandler,
			true,
			true,
		},
		Route{
			"Challonge",
			[]string{"!challonge"},
			challonge.ChallongeHandler(config.ChallongeApi, config.Subdomain, []string{config.AdminChannel, config.PostChannel}, []string{config.AdminChannel}),
			true,
			true,
		},
		Route{
			"Invite",
			[]string{"!invite"},
			inviteHandler,
			false,
			true,
		},
		Route{
			"Twitter",
			[]string{"!twitter", "!tweet", "!announce"},
			twitterHandler,
			true,
			true,
		},
		Route{
			"Make Role",
			[]string{"!makerole"},
			makeRoleHandler,
			true,
			true,
		},
		Route{
			"Remove Role",
			[]string{"!removerole"},
			removeRoleHandler,
			true,
			true,
		},
		Route{
			"Show Roles",
			[]string{"!showroles"},
			showRolesHandler,
			true,
			true,
		},
		Route{
			"Give Role",
			[]string{"!iam "},
			iamHandler,
			false,
			true,
		},
		Route{
			"Take Role",
			[]string{"!iamn"},
			iamnHandler,
			false,
			true,
		},
		Route{
			"Role Call",
			[]string{""},
			RoleCallHandler,
			false,
			false,
		},
		Route{
			"Help Handler",
			[]string{"!help"},
			RoleCallHandler,
			false,
			true,
		},
		Route{
			"Chloe hi",
			[]string{"!hi"},
			hiHandler,
			false,
			true,
		},
		Route{
			"Lanes",
			[]string{"!lanes"},
			lanesHandler,
			true,
			true,
		},
	}
}

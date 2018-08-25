package server

import (
	"github.com/bwmarrin/discordgo"
	"github.com/heshoots/discordbot/pkg/challonge"
)

type HandlerFunc func(s *discordgo.Session, m *discordgo.MessageCreate)

type Route struct {
	Name     string
	Prefix   []string
	Handler  HandlerFunc
	Admin    bool
	HelpText string
}

type Routes []Route

func GetRoutes() Routes {
	return Routes{
		Route{
			"Announce",
			[]string{"!discord", "!announce"},
			discordHandler,
			true,
			"(message) Send Message to announcements channel",
		},
		Route{
			"Challonge",
			[]string{"!challonge"},
			challonge.ChallongeHandler(config.ChallongeApi, config.Subdomain, []string{config.AdminChannel, config.PostChannel}, []string{config.AdminChannel}),
			true,
			"(tournament_name, game_name) Starts challonge tournament with name for the game",
		},
		Route{
			"Invite",
			[]string{"!invite"},
			inviteHandler,
			false,
			"Get invite link",
		},
		Route{
			"Twitter",
			[]string{"!twitter", "!tweet", "!announce"},
			twitterHandler,
			true,
			"Send message to twitter",
		},
		Route{
			"Make Role",
			[]string{"!makerole"},
			makeRoleHandler,
			true,
			"(role_name) Enables role to be added from role_channel",
		},
		Route{
			"Remove Role",
			[]string{"!removerole"},
			removeRoleHandler,
			true,
			"(role_name) Disable role addition from role_channel",
		},
		Route{
			"Show Roles",
			[]string{"!showroles"},
			showRolesHandler,
			true,
			"Display roles in role channel",
		},
		Route{
			"Give Role",
			[]string{"!iam "},
			iamHandler,
			false,
			"(role_name) Get a role",
		},
		Route{
			"Take Role",
			[]string{"!iamn"},
			iamnHandler,
			false,
			"(role_name) Remove a role",
		},
		Route{
			"Help Handler",
			[]string{"!help"},
			helpHandler,
			false,
			"Get help",
		},
		Route{
			"Chloe hi",
			[]string{"!hi"},
			hiHandler,
			false,
			"Say Hello",
		},
		Route{
			"Lanes",
			[]string{"!lanes"},
			lanesHandler,
			true,
			"Check lanes events",
		},
	}
}

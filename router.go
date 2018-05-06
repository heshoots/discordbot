package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func NewRouter(apiKey string) (*discordgo.Session, error) {
	var discord *discordgo.Session
	discord, err := discordgo.New("Bot " + apiKey)
	if err != nil {
		log.Panic("Error creating Discord session", err)
		return nil, err
	}
	routes := GetRoutes()
	for _, route := range routes {
		for _, prefix := range route.Prefix {
			var handler func(s *discordgo.Session, m *discordgo.MessageCreate)
			if route.Logged {
				handler = Logger(route)
			} else {
				handler = route.Handler
			}
			if route.Admin {
				handler = isAdminHandler(handler)
			}
			handler = prefixHandler(prefix, handler)
			discord.AddHandler(handler)
		}
	}
	return discord, nil
}

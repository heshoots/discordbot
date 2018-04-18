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
	for _, route := range routes {
		for _, prefix := range route.Prefix {
			var handler func(s *discordgo.Session, m *discordgo.MessageCreate)
			handler = route.Handler
			handler = prefixHandler(prefix, handler)
			discord.AddHandler(handler)
		}
	}
	return discord, nil
}

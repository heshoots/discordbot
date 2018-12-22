package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/heshoots/discordbot/pkg/models"
	"github.com/heshoots/discordbot/pkg/server"
)

var compiled string

func main() {
	server.SetConfig()
	config := server.GetConfig()
	models.RoleConfig()
	//models.LoadRoles("./config/config.yaml")
	discord, err := server.NewRouter(config.DiscordApi)
	if err != nil {
		log.Panic(err)
		return
	}
	err = discord.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}
	discord.ChannelMessageSend(config.AdminChannel, "Redeployed, compiled: "+compiled)
	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

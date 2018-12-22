package main

import (
	"fmt"
	"log"
	"os"

	"github.com/heshoots/discordbot/pkg/server"
)

var compiled string

func main() {
	server.SetConfig()
	config := server.GetConfig()
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
	fmt.Printf(os.Args[1])
	r, _ := discord.GuildRoles(os.Args[1])
	for _, i := range r {
		fmt.Printf("  - name: %s\n    roleID: %s\n", i.Name, i.ID)
	}
	discord.Close()
}

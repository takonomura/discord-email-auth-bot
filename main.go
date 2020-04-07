package main

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("loading config: %+v", err)
	}
	users, err := LoadUsers()
	if err != nil {
		log.Fatalf("loading users: %+v", err)
	}

	db, err := NewDB("file:authentications.db")
	if err != nil {
		log.Fatalf("creating db: %+v", err)
	}
	err = db.InitializeTables()
	if err != nil {
		log.Fatalf("initializing db: %+v", err)
	}

	sg := &SendGrid{
		Client:    sendgrid.NewSendClient(cfg.SendGridAPIKey),
		Sender:    cfg.Sender,
		GuildName: cfg.GuildName,
		BaseURL:   cfg.BaseURL,
	}

	bot := &Bot{
		Users:    users,
		DB:       db,
		SendGrid: sg,
	}
	discord, err := bot.Start(cfg.DiscordToken)
	if err != nil {
		log.Fatalf("starting bot %+v", err)
	}

	server := &HTTPServer{
		Discord:        discord,
		DiscordGuildID: cfg.DiscordGuild,
		DiscordRoles:   cfg.DiscordRoles,
		Users:          users,
		DB:             db,
	}
	log.Fatal(server.ListenAndServe(":8080"))
}

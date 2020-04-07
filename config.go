package main

import (
	"encoding/json"
	"os"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Config struct {
	DiscordToken string            `json:"discord_token"`
	DiscordGuild string            `json:"discord_guild"`
	DiscordRoles map[string]string `json:"discord_roles"`

	SendGridAPIKey string `json:"sendgrid_api_key"`

	Sender *mail.Email `json:"sender"`

	BaseURL   string `json:"base_url"`
	GuildName string `json:"guild_name"`
}

func LoadConfig() (*Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type User struct {
	ID     string
	Email  string
	Groups []string
}

func LoadUsers() ([]User, error) {
	f, err := os.Open("users.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s []User
	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func FindUser(s []User, id string) *User {
	for _, u := range s {
		if u.ID == id {
			return &u
		}
	}
	return nil
}

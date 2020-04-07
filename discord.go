package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Users    []User
	DB       *DB
	SendGrid *SendGrid
}

func (bot *Bot) Start(token string) (*discordgo.Session, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	s.AddHandler(bot.onMessage)
	err = s.Open()
	return s, err
}

func (bot *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if len(m.Mentions) != 1 || m.Mentions[0].ID != s.State.User.ID {
		return
	}
	log.Print("Received message:" + m.Content)
	content := m.Content
	content = strings.ReplaceAll(content, "<@"+s.State.User.ID+">", "")
	content = strings.ReplaceAll(content, "<@!"+s.State.User.ID+">", "")
	content = strings.TrimSpace(content)
	u := FindUser(bot.Users, content)
	if u == nil {
		_, err := s.ChannelMessageSend(m.ChannelID, "`"+content+"` に該当するデータが見つかりませんでした")
		if err != nil {
			log.Printf("sending error message: %+v", err)
		}
		return
	}

	a, err := bot.DB.StartAuth(u.ID, m.Author.ID)
	if err != nil {
		log.Printf("starting auth: %+v", err)
		_, err := s.ChannelMessageSend(m.ChannelID, "エラーが発生しました")
		if err != nil {
			log.Printf("sending error message: %+v", err)
		}
		return
	}

	err = bot.SendGrid.Send(u.Email, m.Author.String(), a.Token)
	if err != nil {
		log.Printf("sending email: %+v", err)
		_, err := s.ChannelMessageSend(m.ChannelID, "エラーが発生しました")
		if err != nil {
			log.Printf("sending error message: %+v", err)
		}
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "認証コードを発行しました")
	if err != nil {
		log.Printf("sending successful message: %+v", err)
	}
}

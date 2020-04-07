package main

import (
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type HTTPServer struct {
	Discord        *discordgo.Session
	DiscordGuildID string
	DiscordRoles   map[string]string
	Users          []User
	DB             *DB
}

func (s *HTTPServer) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/confirm", s.handleConfirm)
	return http.ListenAndServe(addr, mux)
}

func (s *HTTPServer) handleConfirm(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(rw, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	token := req.URL.Query().Get("token")
	if token == "" {
		http.Error(rw, "400 Bad Request", http.StatusBadRequest)
		return
	}

	log.Printf("Confirming %q", token)

	a, err := s.DB.ConfirmToken(token)
	if err == ErrTokenNotExists {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}
	if err == ErrAlreadyAuthenticated {
		http.Error(rw, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("confirming token: %+v", err)
		http.Error(rw, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	u := FindUser(s.Users, a.UserID)
	if u == nil {
		log.Printf("no user data for %q", a.UserID)
		http.Error(rw, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	var hasError bool
	for _, g := range u.Groups {
		role := s.DiscordRoles[g]
		if role == "" {
			log.Printf("no role for %q", g)
			hasError = true
			continue
		}
		err := s.Discord.GuildMemberRoleAdd(s.DiscordGuildID, a.DiscordID, role)
		if err != nil {
			log.Printf("adding role member: %+v", err)
			hasError = true
		}
	}
	if hasError {
		http.Error(rw, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Error(rw, "Authentication Completed", http.StatusOK)
}

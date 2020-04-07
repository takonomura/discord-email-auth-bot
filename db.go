package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrTokenNotExists       = errors.New("token not exists")
	ErrAlreadyAuthenticated = errors.New("already authenticated")
)

type DB struct {
	DB *sql.DB
}

func NewDB(dsn string) (*DB, error) {
	sqlite, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &DB{DB: sqlite}, nil
}

type Authentication struct {
	Token          string
	TokenExpiredAt time.Time

	UserID    string
	DiscordID string

	CreatedAt       time.Time
	AuthenticatedAt *time.Time
}

func (db *DB) InitializeTables() error {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS authentications (
		token VARCHAR(255) PRIMARY KEY,
		token_expired_at DATETIME NOT NULL,
		user_id VARCHAR(255) NOT NULL,
		discord_id VARCHAR(255) NOT NULL,
		created_at DATETIME NOT NULL,
		authenticated_at DATETIME
	)`)
	return err
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (db *DB) StartAuth(userID, discordID string) (*Authentication, error) {
	t, err := generateToken()
	if err != nil {
		return nil, err
	}
	a := Authentication{
		Token:          t,
		TokenExpiredAt: time.Now().Add(30 * time.Minute),
		UserID:         userID,
		DiscordID:      discordID,
		CreatedAt:      time.Now(),
	}
	_, err = db.DB.Exec(`INSERT INTO authentications (token, token_expired_at, user_id, discord_id, created_at) VALUES (?, ?, ?, ?, ?)`, a.Token, a.TokenExpiredAt, a.UserID, a.DiscordID, a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (db *DB) ConfirmToken(token string) (*Authentication, error) {
	a := Authentication{
		Token: token,
	}
	err := db.DB.QueryRow("SELECT token_expired_at, user_id, discord_id, authenticated_at FROM authentications WHERE token = ?", token).Scan(&a.TokenExpiredAt, &a.UserID, &a.DiscordID, &a.AuthenticatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrTokenNotExists
	}
	if err != nil {
		return nil, err
	}
	if a.TokenExpiredAt.Before(time.Now()) {
		return nil, ErrTokenNotExists
	}
	if a.AuthenticatedAt != nil {
		return &a, ErrAlreadyAuthenticated
	}
	now := time.Now()
	a.AuthenticatedAt = &now
	_, err = db.DB.Exec("UPDATE authentications SET authenticated_at = ? WHERE token = ?", a.AuthenticatedAt, a.Token)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

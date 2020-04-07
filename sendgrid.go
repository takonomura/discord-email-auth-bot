package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct {
	Client *sendgrid.Client

	Sender *mail.Email

	GuildName string
	BaseURL   string
}

func (sg *SendGrid) buildContent(discordTag, token string) string {
	return strings.Join([]string{
		fmt.Sprintf("%s の %s です。", sg.GuildName, sg.Sender.Name),
		"以下の Discord タグが自身のものであることを確認し、問題が無い場合はリンクをクリックして認証を完了してください。",
		"",
		fmt.Sprintf("Discord Tag: %s", discordTag),
		"",
		fmt.Sprintf("%sconfirm?token=%s", sg.BaseURL, token),
		"",
		"このメールに見覚えがない場合は、そのまま無視してください。",
	}, "\n")
}

func (sg *SendGrid) Send(toEmail, discordTag, token string) error {
	subject := sg.GuildName + " の認証"
	to := mail.NewEmail(discordTag, toEmail)
	content := mail.NewContent("text/plain", sg.buildContent(discordTag, token))

	m := mail.NewV3MailInit(sg.Sender, subject, to, content)
	resp, err := sg.Client.Send(m)
	if err != nil {
		return err
	}

	log.Printf("Sent a email to %q", toEmail)
	log.Printf("  Status Code: %d", resp.StatusCode)
	log.Printf("  Body: %q", resp.Body)
	return nil
}

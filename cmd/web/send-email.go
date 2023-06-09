package main

import (
	"log"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for msg := range app.MailChan {
			sendMsg(&msg)
		}
	}()
}

func sendMsg(m *models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	// In production
	// server.Username
	// server.Password
	// server.Encryption

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)
	if err := email.Send(client); err != nil {
		log.Println(err)
	} else {
		log.Println("email sent")
	}
}

package main

import (
	"log"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

// listenForMail is a goroutine that listens for mail messages from the app.MailChan channel
// and sends them using the sendMsg function.
func listenForMail() {
	go func() {
		for msg := range app.MailChan {
			sendMsg(&msg)
		}
	}()
}

// sendMsg sends an email using the provided mail data.
// It connects to the SMTP server and sends the email using the go-simple-mail library.
func sendMsg(m *models.MailData) {
	// Create a new SMTP client
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

	// Connect to the SMTP server
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	// Create a new email message
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	// Send the email
	if err := email.Send(client); err != nil {
		log.Println(err)
	} else {
		log.Println("email sent")
	}
}

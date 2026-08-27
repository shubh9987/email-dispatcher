package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"sync"
	"time"
)

func emailWorker(id int, ch chan Recipient, wg *sync.WaitGroup) {
	defer wg.Done()

	// Read SMTP configuration from .env
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Create SMTP authentication
	auth := smtp.PlainAuth(
		"",
		username,
		password,
		smtpHost,
	)

	for recipient := range ch {

		// Generate email body from template
		msg, err := executeTemplate(recipient)
		if err != nil {
			fmt.Printf("Worker %d: Error parsing template for %s\n", id, recipient.Email)
			continue
		}

		fmt.Printf("Worker %d: Sending email to %s\n", id, recipient.Email)

		// Send email
		err = smtp.SendMail(
			smtpHost+":"+smtpPort,
			auth,
			username,
			[]string{recipient.Email},
			[]byte(msg),
		)

		if err != nil {
			log.Printf("Worker %d: Failed to send email to %s: %v\n", id, recipient.Email, err)
			continue
		}

		fmt.Printf("Worker %d: Successfully sent email to %s\n", id, recipient.Email)

		// Small delay between emails
		time.Sleep(50 * time.Millisecond)
	}
}

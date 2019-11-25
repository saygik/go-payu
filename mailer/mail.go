package mailer

import (
	"log"
	"net/smtp"
)

func Send3(body string) error {
	from := "crystalrentalcarsender@gmail.com"
	pass := "qwe123QWE@"
	to := "crystalrentalcarsender@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print("sent, visit crystalrentalcarsender@gmail.com")
	return nil
}

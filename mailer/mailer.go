package mailer

import (
	"bytes"
	"fmt"
	"github.com/nuveo/log"
	"html/template"
	"net/smtp"
)

// https://hackernoon.com/sending-html-email-using-go-c464d03a26a6
type MailRequest struct {
	from     string
	to       []string
	subject  string
	body     string
	password string
	user     string
	server   string
	port     int
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewMailRequest(user string, password string, server string, port int, to []string, subject string) *MailRequest {

	return &MailRequest{
		to:       to,
		subject:  subject,
		user:     user,
		password: password,
		server:   server,
		port:     port,
	}
}

func (r *MailRequest) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *MailRequest) sendMail() bool {

	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%d", r.server, r.port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", r.user, r.password, r.server), r.user, r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

func (r *MailRequest) Send(templateName string, items interface{}) bool {

	err := r.parseTemplate(templateName, items)
	if err != nil {
		return false
		log.Errorln(err)
	}
	if ok := r.sendMail(); ok {
		log.Println("Email has been sent to %s\n", r.to)
		return true
	} else {
		log.Errorln("Failed to send the email to %s\n", r.to)
		return false
	}
}

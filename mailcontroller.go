package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-payu/mailer"
)

func sendMail(c *gin.Context) {
	fmt.Println("---------sendin eMail")

	subject := "Информация о бронировании автомобиля"
	mr := mailer.NewMailRequest(AppConfig.Mail.User,
		AppConfig.Mail.Password,
		AppConfig.Mail.Server,
		AppConfig.Mail.Port,
		[]string{AppConfig.Mail.MailTo},
		subject)
	mailResult := mr.Send("config/mailtemplate.html",
		map[string]string{
			"username": "Conor", "paymentId": "7W76JLH6XQ191124GUEST000P01"})

	if !mailResult {
		c.JSON(400, gin.H{"Message": "mail not send"})
	} else {
		c.JSON(200, gin.H{"Message": "sended"})
	}
}

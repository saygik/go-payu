package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-payu/mailer"
)

type MailCarsStruct struct {
	Mailto string         `json:"mailto,omitempty"`
	Subj   string         `json:"subj,omitempty"`
	Images []ImagesStruct `json:"images,omitempty"`
}

type ImagesStruct struct {
	Id  string `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
}

func sendMail(c *gin.Context) {
	fmt.Println("---------sendin eMail")
	var img = []string{"", "", "", "", "", ""}
	var MailCars MailCarsStruct
	c.BindJSON(&MailCars)

	for key, _ := range MailCars.Images {
		img[key] = MailCars.Images[key].Url
	}
	fmt.Println(AppConfig.MailCars.User)
	mr := mailer.NewMailRequest(AppConfig.MailCars.User,
		AppConfig.MailCars.Password,
		AppConfig.MailCars.Server,
		AppConfig.MailCars.Port,
		[]string{MailCars.Mailto},
		MailCars.Subj)
	mailResult := mr.Send("templates/mailtemplateCars.html",
		map[string]string{
			"username": "Conor", "img1": img[0],
			"img2": img[1],
			"img3": img[2]})

	if !mailResult {
		c.JSON(400, gin.H{"Message": "mail not send"})
	} else {
		c.JSON(200, gin.H{"Message": "sended"})
	}
}

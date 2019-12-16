package main

import (
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-payu/mailer"
)

type MailCarsStruct struct {
	Mailto               string         `json:"mailto,omitempty"`
	Subj                 string         `json:"subj,omitempty"`
	Name                 string         `json:"name,omitempty"`
	Phone                string         `json:"phone,omitempty"`
	Email                string         `json:"email,omitempty"`
	Brend                string         `json:"brend,omitempty"`
	Model                string         `json:"model,omitempty"`
	Year                 string         `json:"year,omitempty"`
	StandortPLZ          string         `json:"standortPLZ,omitempty"`
	Kilometerstand       string         `json:"kilometerstand,omitempty"`
	Kraftstoffart        string         `json:"kraftstoffart,omitempty"`
	Hubraum              string         `json:"hubraum,omitempty"`
	Getriebe             string         `json:"getriebe,omitempty"`
	MangelSchaden        string         `json:"mangelSchaden,omitempty"`
	IhrePreisvorstellung string         `json:"ihrePreisvorstellung,omitempty"`
	Images               []ImagesStruct `json:"images,omitempty"`
}

type ImagesStruct struct {
	Id  string `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
}

func sendMail(c *gin.Context) {
	//	fmt.Println("---------sendin eMail")
	var img = []string{"", "", "", "", "", ""}
	var imgTitle = []string{"", "", "", "", "", ""}
	var MailCars MailCarsStruct
	c.BindJSON(&MailCars)
	//<tr class="footer"><td style="padding: 40px;">
	//<a href="{{ .img1 }}" target="_blank">
	//<img alt="foto 1" src="{{ .img1 }}"  class="images" />
	//{{ .img1 }}
	//</a></td></tr>
	for key, _ := range MailCars.Images {
		img[key] = MailCars.Images[key].Url
		imgTitle[key] = "Foto" + MailCars.Images[key].Id
	}
	mr := mailer.NewMailRequest(AppConfig.MailCars.User,
		AppConfig.MailCars.Password,
		AppConfig.MailCars.Server,
		AppConfig.MailCars.Port,
		[]string{MailCars.Mailto},
		MailCars.Subj)
	mailResult := mr.Send("templates/mailtemplateCars.html",
		map[string]string{
			"Name":                 MailCars.Name,
			"Phone":                MailCars.Phone,
			"Email":                MailCars.Email,
			"Brend":                MailCars.Brend,
			"Model":                MailCars.Model,
			"Year":                 MailCars.Year,
			"StandortPLZ":          MailCars.StandortPLZ,
			"Kilometerstand":       MailCars.Kilometerstand,
			"Kraftstoffart":        MailCars.Kraftstoffart,
			"Hubraum":              MailCars.Hubraum,
			"Getriebe":             MailCars.Getriebe,
			"MangelSchaden":        MailCars.MangelSchaden,
			"IhrePreisvorstellung": MailCars.IhrePreisvorstellung,
			"img1":                 img[0],
			"img1Title":            imgTitle[0],
			"img2":                 img[1],
			"img2Title":            imgTitle[1],
			"img3":                 img[2],
			"img3Title":            imgTitle[2],
			"img4":                 img[3],
			"img4Title":            imgTitle[3],
			"img5":                 img[4],
			"img5Title":            imgTitle[4],
			"img6":                 img[5],
			"img6Title":            imgTitle[5]})

	if !mailResult {
		c.JSON(400, gin.H{"Message": "mail not send"})
	} else {
		c.JSON(200, gin.H{"Message": "sended"})
	}
}

package mailer

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/smtp"
	"net/url"

	"addysnip.dev/api/pkg/database"
	"addysnip.dev/api/pkg/logger"
	"addysnip.dev/api/pkg/utils"
	models "addysnip.dev/types/database"
)

type Payload struct {
	Template  string                 `json:"template"`
	Data      map[string]interface{} `json:"data"`
	Recipient string                 `json:"recipient"`
	From      string                 `json:"from"`
}

var log = logger.Category("pkg/mailer")

func Handle(body string) error {
	payload := Payload{}
	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		return err
	}

	tem, err := FindTemplate(payload.Template)
	if err != nil {
		return err
	}

	out, err := BuildTemplate(*tem, tem.Name, payload.Data)
	if err != nil {
		return err
	}

	err = Send(payload.Recipient, tem.Subject, BuildBody(payload.From, payload.Recipient, tem.Subject, out))
	if err != nil {
		return err
	}

	return nil
}

func urlEscape(text string) string {
	return url.QueryEscape(text)
}

func FindTemplate(name string) (*models.Template, error) {
	tem := models.Template{}
	if err := database.DB.Where("name = ?", name).First(&tem).Error; err != nil {
		log.Error("Error finding template %s: %s", name, err)
		return nil, err
	}
	return &tem, nil
}

func BuildTemplate(tmpl models.Template, name string, data map[string]interface{}) (*bytes.Buffer, error) {
	t, err := template.New(name).Funcs(template.FuncMap{
		"urlEscape": urlEscape,
	}).Parse(tmpl.Body)
	if err != nil {
		log.Error("Error parsing template %s: %s", name, err)
		return nil, err
	}
	out := new(bytes.Buffer)
	err = t.Execute(out, data)
	if err != nil {
		log.Error("Error executing template %s: %s", name, err)
		return nil, err
	}

	return out, nil
}

func BuildBody(from string, to string, subject string, body *bytes.Buffer) []byte {
	var msg string
	msg = "To: " + to + "\r\n"
	if from != "" {
		msg += "From: " + from + "\r\n"
	} else {
		msg += "From: AddySnip <no-reply@addysnip.com>\r\n"
	}
	msg += "Subject: " + subject + "\r\n"
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	msg += "\r\n" + body.String()

	return []byte(msg)
}

func Send(to string, subject string, body []byte) error {
	log.Debug("Sending email to %s", to)
	auth := smtp.PlainAuth("", utils.Getenv("SMTP_USER", "no-reply@addysnip.com"), utils.Getenv("SMTP_PASS", "password"), utils.Getenv("SMTP_HOST", "mail.addysnip.com"))

	err := smtp.SendMail(utils.Getenv("SMTP_HOST", "mail.addysnip.com")+":"+utils.Getenv("SMTP_PORT", "587"), auth, utils.Getenv("SMTP_USER", "no-reply@addysnip.com"), []string{to}, body)
	if err != nil {
		log.Error("Error sending email: %s", err)
		return err
	}

	log.Debug("Email to %s sent.", to)
	return nil
}

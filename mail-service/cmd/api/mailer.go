package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	smtp "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (mail *Mail) SendSMTPMessage(msg Message) error {

	if msg.From == "" {
		msg.From = mail.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = mail.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := mail.buildHTMLMessage(msg)

	if err != nil {
		return err
	}

	plainTextMessage, err := mail.buildPlainTextMessage(msg)

	if err != nil {
		return err
	}

	server := smtp.NewSMTPClient()
	server.Host = mail.Host
	server.Port = mail.Port
	server.Username = mail.Username
	server.Password = mail.Password
	server.Encryption = mail.getEncyption(mail.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	email := smtp.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)
	email.SetBody(smtp.TextPlain, plainTextMessage)
	email.AddAlternative(smtp.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}
	err = email.Send(smtpClient)

	if err != nil {
		return err
	}
	return nil
}

// build html version of email
func (mail *Mail) buildHTMLMessage(msg Message) (string, error) {

	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-template").ParseFiles(templateToRender)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = mail.inlineCss(formattedMessage)

	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// apply css for a message
func (mail *Mail) inlineCss(formattedMessage string) (string, error) {

	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(formattedMessage, &options)

	if err != nil {
		return "", err
	}

	html, err := prem.Transform()

	if err != nil {
		return "", err
	}
	return html, nil
}

// build plain text message
func (mail *Mail) buildPlainTextMessage(msg Message) (string, error) {

	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainTextMessage := tpl.String()

	if err != nil {
		return "", err
	}

	return plainTextMessage, nil
}

// wich encryption to use for sending message
func (mail *Mail) getEncyption(encry string) smtp.Encryption {

	switch encry {
	case "tls":
		return smtp.EncryptionSTARTTLS
	case "ssl":
		return smtp.EncryptionSSLTLS
	case "none":
		return smtp.EncryptionNone
	default:
		return smtp.EncryptionSTARTTLS
	}
}

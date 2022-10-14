package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type Msg struct {
	Subject    string
	Body       string
	BodyType   string
	Attachment string
}

var EMAIL_TOKEN = "this_is_email_token"
var defaultMsg = &Msg{
	Subject:    "",
	Body:       "",
	BodyType:   "text/html",
	Attachment: "",
}

func InitMsg(subject, body, attachment string) {
	defaultMsg.Subject = subject
	defaultMsg.Body = body
	defaultMsg.Attachment = attachment
}

func SendMail(sender string, receiver string, token string) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", sender)
	mail.SetHeader("To", receiver)
	mail.SetHeader("Subject", defaultMsg.Subject)
	mail.SetBody(defaultMsg.BodyType, defaultMsg.Body)
	if len(defaultMsg.Attachment) != 0 {
		mail.Attach(defaultMsg.Attachment)
	}

	dial := gomail.NewDialer("smtp.gmail.com", 587, sender, token)
	if err := dial.DialAndSend(mail); err != nil {
		fmt.Printf("err %v", dial)
		return
	}
}

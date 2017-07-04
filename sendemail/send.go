package sendemail

import (
	"log"
	"net/smtp"
	"strconv"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}
type SmtpAuth struct {
	user EmailUser
	auth smtp.Auth
}

// AuthMailBox authorises your email server to send notifications
func AuthMailBox(user EmailUser) (auth SmtpAuth) {
	auth1 := smtp.PlainAuth("",
		user.Username,
		user.Password,
		user.EmailServer,
	)
	a := SmtpAuth{
		user,
		auth1,
	}
	return a
}

//SendEmailwithMessage sends message to user with message about metric out of boundaries
func SendEmailwithMessage(addres, msg string, auth SmtpAuth) error {

	var err error
	log.Printf("%+v", auth, "Email Auth Info")
	msg2 := []byte(msg)

	err = smtp.SendMail(auth.user.EmailServer+":"+strconv.Itoa(auth.user.Port),
		auth.auth,
		auth.user.Username,
		[]string{addres},
		msg2)
	if err != nil {
		log.Print("SendEmailwithMessage ERROR: attempting to send a mail ", err)
		return err
	}
	return nil
}

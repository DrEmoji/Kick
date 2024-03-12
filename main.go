package main

import (
	Kick "app/extra/kick"
	"fmt"
	"strings"
	"time"
)

func main() {
	email := Kick.GenEmail()
	kpsdkct := Kick.GetCT()
	kick := Kick.CreateClient(Kick.GenUsername("GhostSpirit"), email.Data.GenerateEmail.Email, "qwswqSws1214>", kpsdkct)
	kick.GetCookies()
	kick.StartSocket()
	defer kick.Conn.Close()
	kick.RequestTokenProvider()
	kick.SendEmail()
	time.Sleep(3 * time.Second)
	mailbox := Kick.GetMailBox(email.Data.GenerateEmail.AccessToken)
	code := strings.Split(mailbox.Data.GetMailList.Mails[0].Subject, " ")[0]
	fmt.Println(code)
	kick.SendEmailCode(code)
	username, err := kick.RegisterAccount()
	if err != nil {
		fmt.Println("an error occured men...")
		fmt.Println(err.Error())
		return
	}
	fmt.Printf(`%s:%s:%s`, username, kick.Email, kick.Password)
	// // kick.RequestTokenProvider()
	kick.LoginAccount()
	time.Sleep(3 * time.Second)
	mailbox = Kick.GetMailBox(email.Data.GenerateEmail.AccessToken)
	fmt.Println(mailbox)
	code = strings.Split(mailbox.Data.GetMailList.Mails[0].Subject, " ")[0]
	kick.SendLoginCode(code)
	kick.SendMessage("6692367", "Hello")
}

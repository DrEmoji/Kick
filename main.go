package main

import (
	Kick "app/extra/kick"
	"fmt"
	"strings"
	"time"
)

func main() {
	email := Kick.GenEmail()
	kpsdkct := "024t5uIBGXiN99hup7i3ivsxbWlFI4T2vNafXkr8gK0VGlZjBfHFi5uOTZnYnu7IVRs7jc5WOVK1bWwtUZKBkh7skN1jPLeRpL6BA4VmzvGlcUpNAUldPBREVaef5hxrXYzsivMkW1R3ALTvgCrEHXP7c36rNH" //Kick.GetCT()
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
	kick.RequestTokenProvider()
	kick.LoginAccount("")
	time.Sleep(3 * time.Second)
	mailbox = Kick.GetMailBox(email.Data.GenerateEmail.AccessToken)
	code = strings.Split(mailbox.Data.GetMailList.Mails[0].Subject, " ")[0]
	kick.LoginAccount(code)
}

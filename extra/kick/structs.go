package kick

import (
	"log"
	"time"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/gorilla/websocket"
)

// Form is to define the response from ""
type Form struct {
	Enabled                   bool   `json:"enabled"`
	NameFieldName             string `json:"nameFieldName"`
	UnrandomizedNameFieldName string `json:"unrandomizedNameFieldName"`
	ValidFromFieldName        string `json:"validFromFieldName"`
	EncryptedValidFrom        string `json:"encryptedValidFrom"`
}

// Client is to define the over-all package of kick
type Client struct {
	Conn     *websocket.Conn
	Username string
	Email    string
	Password string
	xKpsdkCt string
	Auth     string
	socketID string
	xsrf     string
	request  tls_client.HttpClient
	form     Form
}

type Email struct {
	Data struct {
		GenerateEmail struct {
			Email       string `json:"email"`
			AccessToken string `json:"accessToken"`
			Typename    string `json:"__typename"`
		} `json:"generateEmail"`
	} `json:"data"`
}

type Mailbox struct {
	Data struct {
		GetMailList struct {
			TotalCount int  `json:"totalCount"`
			HasMore    bool `json:"hasMore"`
			Mails      []struct {
				ID        string    `json:"id"`
				Name      string    `json:"name"`
				Address   string    `json:"address"`
				Subject   string    `json:"subject"`
				CreatedAt time.Time `json:"createdAt"`
				Typename  string    `json:"__typename"`
			} `json:"mails"`
			Typename string `json:"__typename"`
		} `json:"getMailList"`
	} `json:"data"`
}

func CreateClient(username, email, password, kpsdkct string) *Client {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Firefox_102),
		tls_client.WithRandomTLSExtensionOrder(),
		tls_client.WithCookieJar(tls_client.NewCookieJar(tls_client.WithLogger(tls_client.NewLogger()))),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Username: username,
		Email:    email,
		xKpsdkCt: kpsdkct,
		Password: password,
		request:  client,
	}
}

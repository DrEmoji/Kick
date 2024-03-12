package kick

import (
	"encoding/json"
	"log"
	"math/rand"

	"net/http"

	"github.com/gorilla/websocket"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

type SocketResponse struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenUsername(prefix string) string {
	return prefix + "_" + RandStringRunes(5)
}

func connectToWebsocket() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial("wss://ws-us2.pusher.com/app/eb1d5f283081a78b932c?protocol=7&client=js&version=7.6.0&flash=false", http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0"},
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return conn
}

func (client *Client) StartSocket() {
	client.Conn = connectToWebsocket()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		var payload SocketResponse
		err = json.Unmarshal(message, &payload)
		if err != nil {
			log.Println(err)
			return
		}

		if payload.Event == "pusher:connection_established" {
			var data map[string]interface{}
			json.Unmarshal([]byte(payload.Data), &data)
			client.socketID = data["socket_id"].(string)
			break
		}
		continue
	}
}

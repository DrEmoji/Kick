package kick

import (
	"bytes"
	"encoding/json"
	"log"

	"net/http"
)

func GenEmail() Email {
	jsonPayload := []byte(`{"operationName":"GenerateEmail","variables":{},"query":"mutation GenerateEmail($name: String, $expiresInSeconds: Int) {\n generateEmail(name: $name, expiresInSeconds: $expiresInSeconds) {\n email\n accessToken\n __typename\n }\n}"}`)

	req, err := http.NewRequest("POST", "https://api.phantom-mail.io/graphql", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Sec-Ch-Ua", "")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5845.141 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", "")
	req.Header.Set("Origin", "https://www.phantom-mail.io/")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.phantom-mail.io/")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var solutionResp Email
	if err := json.NewDecoder(resp.Body).Decode(&solutionResp); err != nil {
		log.Fatal(err)
	}
	return solutionResp
}

func GetMailBox(auth string) Mailbox {
	jsonPayload := []byte(`{"operationName":"GetMailList","variables":{"take":5,"skip":0},"query":"query GetMailList($skip: Int, $take: Int) {\n  getMailList(skip: $skip, take: $take) {\n    totalCount\n    hasMore\n    mails {\n      id\n      name\n      address\n      subject\n      createdAt\n      __typename\n    }\n    __typename\n  }\n}"}`)

	req, err := http.NewRequest("POST", "https://api.phantom-mail.io/graphql", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Sec-Ch-Ua", "")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5845.141 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", "")
	req.Header.Set("Authorization", "Bearer "+auth)
	req.Header.Set("Origin", "https://www.phantom-mail.io/")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.phantom-mail.io/")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var solutionResp Mailbox
	if err := json.NewDecoder(resp.Body).Decode(&solutionResp); err != nil {
		log.Fatal(err)
	}
	return solutionResp
}

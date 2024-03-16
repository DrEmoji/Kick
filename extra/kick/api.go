package kick

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (client *Client) grabCookie(domain, cookieName string) string {
	_domain, err := url.Parse(fmt.Sprintf(`https://%s`, domain))
	if err != nil {
		log.Fatal(err)
	}

	for _, cookie := range client.request.GetCookies(_domain) {
		if strings.EqualFold(strings.ToLower(cookie.Name), strings.ToLower(cookieName)) {
			return cookie.Value
		}
	}

	return ""
}

// func (client *Client) printCookies(domain string) string {
// 	_domain, err := url.Parse(fmt.Sprintf(`https://%s`, domain))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, cookie := range client.request.GetCookies(_domain) {
// 		fmt.Println("Name:", cookie.Name)
// 		fmt.Println("Value:", cookie.Value)
// 		fmt.Println("Domain:", cookie.Domain)
// 		fmt.Println("Path:", cookie.Path)
// 		fmt.Println("Expires:", cookie.Expires)
// 		fmt.Println("Secure:", cookie.Secure)
// 		fmt.Println("HttpOnly:", cookie.HttpOnly)
// 		fmt.Println()
// 	}

// 	return ""
// }

func GetCD() string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:5000/cd_solution", nil)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("error: on requesting solution..")
		return ""
	}

	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		log.Fatal(err)
		return ""
	}

	xKpsdkCd, ok := jsonResponse["x-kpsdk-cd"].(string)
	if !ok {
		log.Println("error: x-kpsdk-cd not found or not a string")
		return ""
	}
	return xKpsdkCd
}

func GetCT() string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:5000/ct_solution", nil)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("error: on requesting solution..")
		return ""
	}

	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		log.Fatal(err)
		return ""
	}

	xKpsdkCt, ok := jsonResponse["x-kpsdk-ct"].(string)
	if !ok {
		log.Println("error: x-kpsdk-cd not found or not a string")
		return ""
	}
	return xKpsdkCt
}

func (client *Client) GetCookies() {
	req, err := http.NewRequest("GET", "https://kick.com/", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header = http.Header{
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language":           {"en-US,en;q=0.5"},
		"DNT":                       {"1"},
		"Upgrade-Insecure-Requests": {"1"},
		"Connection":                {"keep-alive"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Pragma":                    {"no-cache"},
		"Cache-Control":             {"no-cache"},
	}
	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println("error: on requesting cookies..")
		return
	}
	log.Println("successfully requested cookies..")
	client.xsrf = client.grabCookie("kick.com", "XSRF-TOKEN")
}

func (client *Client) RequestTokenProvider() {
	log.Println(client.socketID)
	req, err := http.NewRequest("GET", "https://kick.com/kick-token-provider", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header = http.Header{
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language":           {"en-US,en;q=0.5"},
		"DNT":                       {"1"},
		"Upgrade-Insecure-Requests": {"1"},
		"Connection":                {"keep-alive"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Pragma":                    {"no-cache"},
		"Cache-Control":             {"no-cache"},
		"X-Socket-ID":               {client.socketID},
		"Referer":                   {"https://kick.com/"},
		"Authorization":             {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":              {client.xsrf},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("error: on requesting cookies..")
		return
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bodyText, &client.form)
	log.Println("successfully requested token provider..")
}

func (client *Client) GetUser() {
	log.Println(client.socketID)
	req, err := http.NewRequest("GET", "https://kick.com/api/v1/user", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header = http.Header{
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language":           {"en-US,en;q=0.5"},
		"DNT":                       {"1"},
		"Upgrade-Insecure-Requests": {"1"},
		"Connection":                {"keep-alive"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Pragma":                    {"no-cache"},
		"Cache-Control":             {"no-cache"},
		"X-Socket-ID":               {client.socketID},
		"Referer":                   {"https://kick.com/"},
		"Authorization":             {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":              {client.xsrf},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("error: on requesting cookies..")
		return
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bodyText))
	json.Unmarshal(bodyText, &client.form)
	log.Println("successfully requested token provider..")
}

func (client *Client) SendEmail() {
	req, err := http.NewRequest("POST", "https://kick.com/api/v1/signup/send/email", strings.NewReader(fmt.Sprintf(`{"email":"%s"}`, client.Email)))
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header = http.Header{
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"},
		"Accept":          {"application/json, text/plain, */*"},
		"Accept-Language": {"en-US"},
		"Content-Type":    {"application/json"},
		"X-Socket-ID":     {client.socketID},
		"Authorization":   {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":    {client.xsrf},
		"Origin":          {"https://kick.com"},
		"DNT":             {"1"},
		"Connection":      {"keep-alive"},
		"X-Kpsdk-Ct":      {client.xKpsdkCt},
		"X-Kpsdk-Cd":      {GetCD()},
		"Referer":         {"https:/`/kick.com/"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-origin"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		fmt.Println(resp.StatusCode)
		log.Println("error: on email send..")
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bodyText))
		return
	}
	log.Println("successfully sent email..")
}

func (client *Client) SendEmailCode(code string) {
	req, err := http.NewRequest("POST", "https://kick.com/api/v1/signup/verify/code", strings.NewReader(fmt.Sprintf(`{"code":"%s","email":"%s"}`, code, client.Email)))
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header = http.Header{
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"},
		"Accept":          {"application/json, text/plain, */*"},
		"Accept-Language": {"en-US"},
		"Content-Type":    {"application/json"},
		"X-Socket-ID":     {client.socketID},
		"Authorization":   {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":    {client.xsrf},
		"Origin":          {"https://kick.com"},
		"DNT":             {"1"},
		"Connection":      {"keep-alive"},
		"Referer":         {"https:/`/kick.com/"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"X-Kpsdk-Ct":      {client.xKpsdkCt},
		"X-Kpsdk-Cd":      {GetCD()},
		"Sec-Fetch-Site":  {"same-origin"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		fmt.Println(resp.StatusCode)
		log.Println("error: on email send..")
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bodyText))
		return
	}

	log.Println("successfully sent email code..")
}

func (client *Client) RegisterAccount() (string, error) {
	payload := fmt.Sprintf(
		`{"birthdate":"03/27/1995","username":"%s","email":"%s","cf_captcha_token":"","password":"%s","enable_sms_promo":false,"enable_sms_security":false,"password_confirmation":"%s","agreed_to_terms":true,"%s":"","_kick_token_valid_from":"%s"}`,
		client.Username,
		client.Email,
		client.Password,
		client.Password,
		client.form.NameFieldName,
		client.form.EncryptedValidFrom,
	)

	req, err := http.NewRequest("POST", "https://kick.com/register", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err)
		return client.Username, err
	}
	req.Header = http.Header{
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"},
		"Accept":          {"application/json, text/plain, */*"},
		"Accept-Language": {"en-US"},
		"Content-Type":    {"application/json"},
		"X-Socket-ID":     {client.socketID},
		"Authorization":   {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":    {client.xsrf},
		"Origin":          {"https://kick.com"},
		"DNT":             {"1"},
		"Connection":      {"keep-alive"},
		"X-Kpsdk-Ct":      {client.xKpsdkCt},
		"X-Kpsdk-Cd":      {GetCD()},
		"Referer":         {"https:/`/kick.com/"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-origin"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
		return client.Username, err
	}
	defer resp.Body.Close()

	// Loop over header names
	for name, values := range resp.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		client.session = client.grabCookie("kick.com", "kick_session")
		fmt.Println(client.session)
		log.Println("successfully sent register..")
		return client.Username, nil
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println("error: on register..")
	fmt.Println(string(bodyText))

	return client.Username, errors.New(string(bodyText))
}

func (client *Client) LoginAccount(code string) {
	payload := fmt.Sprintf(
		`{"email":"%s","password":"%s","_kick_token_IRw3oA2PPsYqA6EB":"","_kick_token_valid_from":"%s","isMobileRequest":true,"one_time_password":"%s"}`,
		client.Email,
		client.Password,
		client.form.EncryptedValidFrom,
		code,
	)

	req, err := http.NewRequest("POST", "https://kick.com/mobile/login", bytes.NewBufferString(payload))
	if err != nil {
		log.Fatal(err)
	}

	req.Header = http.Header{
		"User-Agent":       {"Mozilla/5.0 (Linux; Android 9; SM-G950F Build/PPR1.180610.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/122.0.6261.119 Mobile Safari/537.36"},
		"Accept":           {"application/json, text/plain, */*"},
		"Accept-Language":  {"en-US"},
		"Content-Type":     {"application/json"},
		"X-Socket-ID":      {client.socketID},
		"Authorization":    {"Bearer " + client.xsrf},
		"X-XSRF-TOKEN":     {client.xsrf},
		"X-App-Platform":   {"Android"},
		"X-Requested-With": {"com.kick.mobile"},
		"Origin":           {"https://kick.com"},
		"DNT":              {"1"},
		"Connection":       {"keep-alive"},
		"Referer":          {"https:/`/kick.com/"},
		"Sec-Fetch-Dest":   {"empty"},
		"Sec-Fetch-Mode":   {"cors"},
		"Sec-Fetch-Site":   {"same-origin"},
		"Pragma":           {"no-cache"},
		"Cache-Control":    {"no-cache"},
	}

	resp, err := client.request.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	var responseMap map[string]interface{}
	errr := json.Unmarshal(body, &responseMap)
	if errr != nil {
		log.Fatal(errr)
	}

	// Check if "token" field exists in the response
	token, tokenExists := responseMap["token"].(string)
	if tokenExists {
		log.Println("Token:", token)
		client.Auth = token
	} else {
		log.Println("Authenticated")
	}
}

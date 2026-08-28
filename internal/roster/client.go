package roster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"shifts-go/internal/helper"
	"time"
)

type Client struct {
	LOGIN_URL  string
	ROSTER_URL string
	EMAIL      string
	PASSWORD   string
}

type LoginRequest struct {
	LoginIp  string `json:"loginIp"`
	Password string `json:"password"`
	UserId   string `json:"userId"`
}

type LoginResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	LoginToken string `json:"loginToken"`
	ExpiryTime int64  `json:"expiryTime"`
}

func NewClient() (*Client, error) {
	loginURL := os.Getenv("LOGIN_URL")
	rosterURL := os.Getenv("ROSTER_URL")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	if loginURL == "" || rosterURL == "" {
		return nil, fmt.Errorf("LOGIN URL AND ROSTER URL be provided in .env")
	}

	if email == "" || password == "" {
		return nil, fmt.Errorf("EMAIL AND PASSWORD must be provided in .env")
	}

	return &Client{LOGIN_URL: loginURL, ROSTER_URL: rosterURL, EMAIL: email, PASSWORD: password}, nil
}

func (client *Client) GetEmployees(time time.Time) (*EmployeeResponse, error) {
	date := helper.FormatTime(time)
	loginToken, err := client.LoginAndGetToken()
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(client.ROSTER_URL)
	if err != nil {
		return nil, err
	}

	q := url.Query()
	q.Set("date", date)
	url.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	headers := http.Header{
		"Content-Type":        {"application/json"},
		"synergy-login-token": {loginToken},
	}

	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, err
	}

	var response EmployeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) LoginAndGetToken() (string, error) {
	cachedToken, found, err := GetCachedToken()
	if err != nil {
		return "", err
	}

	if found && len(cachedToken) > 0 {
		return cachedToken, nil
	}

	headers := http.Header{
		"Content-Type": {"application/json"},
	}

	body := LoginRequest{
		LoginIp:  "",
		Password: os.Getenv("PASSWORD"),
		UserId:   os.Getenv("USER_ID"),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		client.LOGIN_URL,
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return "", err
	}

	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("login failed: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response LoginResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}

	cache := CacheResponse{
		LoginToken: response.LoginToken,
		ExpiryDate: response.ExpiryTime,
	}

	if err := Write("token.json", cache); err != nil {
		return "", err
	}

	return response.LoginToken, nil
}

func GetCachedToken() (string, bool, error) {
	if !IsExist("token.json") {
		return "", false, nil
	}

	if IsEmpty("token.json") {
		return "", false, nil
	}

	contents, success, err := Read("token.json")
	if err != nil {
		return "", false, err
	}

	if !success {
		return "", false, nil
	}

	var response CacheResponse
	if err := json.Unmarshal(contents, &response); err != nil {
		return "", false, err
	}

	if time.Now().Unix() >= response.ExpiryDate {
		return "", false, nil
	}

	return response.LoginToken, true, nil
}

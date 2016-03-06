package iotkit

import (
	"net/http"

	"time"
)

// https://github.com/enableiot/iotkit-api/wiki/Authorization

type AuthorizationApi struct {
	client *Client
}

type Account struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type TokenInfo struct {
	Header struct {
		Typ string `json:"typ"`
		Alg string `json:"alg"`
	} `json:"header"`
	Payload struct {
		Jti      string    `json:"jti"`
		Iss      string    `json:"iss"`
		Sub      string    `json:"sub"`
		Exp      time.Time `json:"exp"`
		Accounts []Account `json:"accounts"`
	} `json:"payload"`
}

func (a *AuthorizationApi) CreateToken(username, password string) (string, *http.Response, error) {
	creds := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{username, password}

	req, err := a.client.NewRequest("POST", "auth/token", creds)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var token struct{ Token string }

	resp, err := a.client.Do(req, &token)
	if err != nil {
		return "", resp, err
	}

	return token.Token, resp, nil
}

// https://github.com/enableiot/iotkit-api/wiki/Authorization#get-user-jwt-token-information
//
func (a *AuthorizationApi) GetTokenInfo(authToken string) (*TokenInfo, *http.Response, error) {

	req, err := a.client.NewRequest("GET", "auth/tokenInfo", nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	var ti = new(TokenInfo)
	resp, err := a.client.Do(req, ti)

	return ti, resp, err
}

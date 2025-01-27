package controllers

import (
	"context"
	"encoding/json"
	"fitness-tracker/env"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Google one-off auth token provided to the server to exchange for the googleTokens pair
type GoogleAuthCodeInput struct {
	Code string `json:"code"`
}

// Authenticate with Google by passing short lived token
func AuthenticateWithGoogle(config env.Config, authCode GoogleAuthCodeInput) (profile *GoogleProfileInfo, err error) {
	// got authcode for Google, request tokens and store them in db
	googleTokens, err := googleOneOffTokenExchange(config, authCode)
	if err != nil {
		log.Printf("main: google token exchange for session failed: %v\n", err)
		return nil, err
	}

	// generate user in db and add Gtokens
	profile, err = googleGetProfileInfo(googleTokens)
	if err != nil {
		log.Printf("main: google profile read failed: %v\n", err)
		return nil, err
	}
	log.Printf("main: customer's info: %v\n", *profile)

	return
}

/**
* Handle one-off token exchange for session/refresh pair
 */
func googleOneOffTokenExchange(config env.Config, authCode GoogleAuthCodeInput) (token *oauth2.Token, err error) {
	var googleOauthConfig = &oauth2.Config{
		// port is of the ReactJS app
		RedirectURL:  config.Env.GoogleClientCallbackUrl(),
		ClientID:     config.Env.GoogleClientId(),
		ClientSecret: config.SecEnv.GoogleClientSecret(),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// Use code to get googleTokens and get user info from Google.
	googleTokens, err := googleOauthConfig.Exchange(context.Background(), authCode.Code)
	if err != nil {
		log.Printf("main: exchange of authcode (%v) failed: %v\n", authCode.Code, err)
		return nil, err
	}

	return googleTokens, nil
}

type GoogleProfileInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func googleGetProfileInfo(googleTokens *oauth2.Token) (profileInfo *GoogleProfileInfo, err error) {
	response, err := http.Get(oauthGoogleUrlAPI + googleTokens.AccessToken)
	if err != nil {
		log.Println("main: getting profile data failed")
		return nil, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("main: reading profile data failed")
		return nil, err
	}

	var googleProfileInfo GoogleProfileInfo
	err = json.Unmarshal(contents, &googleProfileInfo)
	if err != nil {
		log.Println("main: reading profile data failed")
		return nil, err
	}
	return &googleProfileInfo, err
}

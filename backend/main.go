package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleClientId string
var googleClientSecret string
var googleClientCallbackUrl string
var serverPort string

var filePathKeyPrivate string
var filePathKeyPublic string

func main() {

	// for now a simple load of credentials
	log.Println("Loading env variables")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!")
	}

	googleClientId = os.Getenv("CLIENT_ID")
	googleClientSecret = os.Getenv("CLIENT_SECRET")
	googleClientCallbackUrl = os.Getenv("CLIENT_CALLBACK_URL")

	filePathKeyPrivate = os.Getenv("FILE_KEY_PRIVATE")
	filePathKeyPublic = os.Getenv("FILE_KEY_PUBLIC")

	if googleClientId == "" || googleClientSecret == "" ||
		googleClientCallbackUrl == "" || filePathKeyPrivate == "" || filePathKeyPublic == "" {
		log.Fatal(`Environment variables (CLIENT_ID, CLIENT_SECRET,
         CLIENT_CALLBACK_URL, FILE_KEY_PRIVATE, FILE_KEY_PUBLIC) are required`)
	}

	serverPort = os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	fmt.Printf("Starting backend on the port: %v", serverPort)
	// starting the server that will listen forever on the port
	http.HandleFunc("/", rootHandler)
	http.Handle("/api/auth/google", corsHandler(googleAuthCodeHandler))

	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}

func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("WEB_HOST"))
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, tokens")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")

		if r.Method == "OPTIONS" {
			return
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

func rootHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if len(request.URL.Path) > 0 {
		responseWriter.WriteHeader(http.StatusNotFound)
		responseWriter.Write([]byte(fmt.Sprintf("URL(%v) is not supported\n", request.URL.Path)))
	} else {
		responseWriter.Write([]byte(fmt.Sprintln("Root endpoint for the backend")))
		responseWriter.Write([]byte(`(debug)Available endpoints:
            \t/api/auth/google
        `))
		log.Println(request.URL.Path)
	}
}

// Google one-off auth token provided to the server to exchange for the tokens pair
type GoogleAuthCodeInput struct {
	Code string `json:"code"`
}

// Return session token in a payload
type SessionTokenOutput struct {
	Token string `json:"token"`
}

func googleAuthCodeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("Reading authCode failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	var authCode GoogleAuthCodeInput
	err = json.Unmarshal(body, &authCode)
	if err != nil {
		log.Printf("Decoding authCode failed: %v; body: %v\n", err, body)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	tokens, err := googleOneOffTokenExchange(authCode)
	if err != nil {
		log.Printf("Google token exchange for session failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: generate JWT token and a refresh token to handle session on the server
	//  and keep it independent from the SSO provider
	// This also allows server to have finer control over data and expiration
	jwt.New(&jwt.SigningMethodRSA{})

	http.SetCookie(responseWriter, &http.Cookie{Name: "gtoken", Value: tokens, Expires: tokens.Expiry})
	json.NewEncoder(responseWriter).Encode(SessionTokenOutput{Token: tokens.AccessToken})
}

func userInfoHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// err := verifyrequest.Header("X-Auth-Token")
	// request.Cookie("gtoken")
}

/**
* Handle one-off token exchange for session/refresh pair
 */
func googleOneOffTokenExchange(authCode GoogleAuthCodeInput) (token *oauth2.Token, err error) {
	var googleOauthConfig = &oauth2.Config{
		// port is of the ReactJS app
		RedirectURL:  fmt.Sprintf("http://localhost:%v", 3000),
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// Use code to get tokens and get user info from Google.
	tokens, err := googleOauthConfig.Exchange(context.Background(), authCode.Code)
	if err != nil {
		log.Printf("Exchange of authcode (%v) failed: %v\n", authCode.Code, err)
		return nil, err
	}

	return tokens, nil
}

type GoogleProfileInfo struct {
	Profile string `json:"profile"`
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func googleGetProfileInfo(tokens *oauth2.Token) (profileInfo *GoogleProfileInfo, err error) {
	response, err := http.Get(oauthGoogleUrlAPI + tokens.AccessToken)
	if err != nil {
		log.Println("Getting profile data failed")
		return nil, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Reading profile data failed")
		return nil, err
	}
	// responseWriter.Header().Add("Content-Type", "application/json")
	// just sending bytes
	log.Printf("Got content:\n %s\n", contents)
	var googleProfileInfo GoogleProfileInfo
	err = json.Unmarshal(contents, &googleProfileInfo)
	if err != nil {
		log.Println("Reading profile data failed")
		return nil, err
	}
	return &googleProfileInfo, err
}

// Create a token
// Ref: https://github.com/golang-jwt/jwt/blob/main/example_test.go
func jwtWithCustomClaims() {
	cryptoKey := LoadECPrivateKeyFromDisk(fmt.Sprintf("./%s", filePathKeyPrivate))

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.RegisteredClaims
	}

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	fmt.Printf("foo: %v\n", claims.Foo)

	// Create claims while leaving out some of the optional fields
	claims = MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)
	ss, err := token.SignedString(cryptoKey)
	fmt.Println(ss, err)

	// Output: foo: bar
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJpc3MiOiJ0ZXN0IiwiZXhwIjoxNTE2MjM5MDIyfQ.xVuY2FZ_MRXMIEgVQ7J-TFtaucVFRXUzHm9LmV41goM <nil>
}

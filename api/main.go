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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const HEADER_JWT = "jwt_token"

var apiPort string

// host excluding port
var cookieDomain string

// host including port for CORS allow
var webBaseUrl string

var googleClientId string
var googleClientSecret string
var googleClientCallbackUrl string

var filePathKeyPrivate string
var filePathKeyPublic string

func main() {
	googleClientId = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleClientCallbackUrl = os.Getenv("GOOGLE_CLIENT_CALLBACK_URL")
	apiPort = os.Getenv("API_PORT")
	cookieDomain = os.Getenv("COOKIE_DOMAIN")
	webBaseUrl = os.Getenv("WEB_BASE_URL")

	filePathKeyPrivate = os.Getenv("FILE_KEY_PRIVATE")
	filePathKeyPublic = os.Getenv("FILE_KEY_PUBLIC")

	if googleClientId == "" || googleClientSecret == "" ||
		googleClientCallbackUrl == "" || filePathKeyPrivate == "" ||
		filePathKeyPublic == "" || cookieDomain == "" || apiPort == "" || webBaseUrl == "" {
		log.Fatalf(`Environment variables 
		(GOOGLE_CLIENT_ID: %v, GOOGLE_CLIENT_SECRET: %v,
         GOOGLE_CLIENT_CALLBACK_URL: %v, 
		 FILE_KEY_PRIVATE: %v, FILE_KEY_PUBLIC: %v,
		 COOKIE_DOMAIN: %v, API_PORT: %v, WEB_BASE_URL: %v) are required`,
			googleClientId,
			"<sensored>",
			googleClientCallbackUrl,
			filePathKeyPrivate,
			filePathKeyPublic,
			cookieDomain,
			apiPort,
			webBaseUrl)
	}

	fmt.Printf("Starting backend on the port: %v", apiPort)
	// starting the server that will listen forever on the port
	http.HandleFunc("/", rootHandler)
	http.Handle("/api/auth/google", corsHandler(googleAuthCodeHandler))
	http.HandleFunc("/authenticated", corsHandler(authenticatedCallHandler))

	log.Fatal(http.ListenAndServe(":"+apiPort, nil))
}

func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", webBaseUrl)
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, googleTokens")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

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

func authenticatedCallHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// testing recovering from panic
	defer func() {
		if err := recover(); err != nil {
			switch t := err.(type) {
			case string:
				http.Error(responseWriter, t, http.StatusUnauthorized)
			default:
				http.Error(responseWriter, "Unspecified panic", http.StatusUnauthorized)
			}
		}
	}()

	jwtCookie, err := request.Cookie(HEADER_JWT)
	if err != nil {
		log.Printf("Token is missing: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	cryptoPublicKey := LoadRSAPublicKeyFromDisk(fmt.Sprintf("./%s", filePathKeyPublic))

	parsed, err := jwt.Parse(jwtCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return cryptoPublicKey, nil
	})

	if !parsed.Valid {
		log.Printf("Token verification failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	// why do this not work while Parse above does?
	// err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), ([]byte)(parts[2]), cryptoPublicKey)
	// if err != nil {
	// 	log.Printf("Token verification failed: %v\n", err)
	// 	http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	a, _ := json.Marshal(map[string]int{"foo": 1, "bar": 2, "baz": 3})
	responseWriter.Write(a)
}

// Google one-off auth token provided to the server to exchange for the googleTokens pair
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

	// got authcode for Google, request tokens and store them in db

	googleTokens, err := googleOneOffTokenExchange(authCode)
	if err != nil {
		log.Printf("Google token exchange for session failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// generate user in db and add Gtokens

	info, err := googleGetProfileInfo(googleTokens)
	if err != nil {
		log.Printf("Google profile read failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Customer's info: %v\n", *info)
	//TODO: add info to DB along with refresh token

	// generating only jwt token for now. Using Google Id as subject in token
	serverJWTToken, err := jwtWithCustomClaims(info.Id, filePathKeyPrivate, time.Now())
	if err != nil {
		log.Printf("JWT server token generated: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// set http-only jwk token cookie
	http.SetCookie(responseWriter, &http.Cookie{
		Name:     HEADER_JWT,
		Value:    serverJWTToken,
		Domain:   cookieDomain,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //NOTE: for testing purposes only
	})
	responseWriter.Write([]byte(""))
}

// func userInfoHandler(responseWriter http.ResponseWriter, request *http.Request) {
// 	// err := verifyrequest.Header("X-Auth-Token")
// 	// request.Cookie("gtoken")
// }

/**
* Handle one-off token exchange for session/refresh pair
 */
func googleOneOffTokenExchange(authCode GoogleAuthCodeInput) (token *oauth2.Token, err error) {
	var googleOauthConfig = &oauth2.Config{
		// port is of the ReactJS app
		RedirectURL:  googleClientCallbackUrl,
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// Use code to get googleTokens and get user info from Google.
	googleTokens, err := googleOauthConfig.Exchange(context.Background(), authCode.Code)
	if err != nil {
		log.Printf("Exchange of authcode (%v) failed: %v\n", authCode.Code, err)
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
		log.Println("Getting profile data failed")
		return nil, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Reading profile data failed")
		return nil, err
	}

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
func jwtWithCustomClaims(customerId, filePathKeyPrivate string, now time.Time) (string, error) {
	cryptoKey := LoadRSAPrivateKeyFromDisk(fmt.Sprintf("./%s", filePathKeyPrivate))

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.RegisteredClaims
	}

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			// For now it's 1 minute for testing purposes
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "FitnessTracker",
			Subject:   customerId,
			Audience:  []string{"FitnessTrackerAPI"},
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)

	return token.SignedString(cryptoKey)
}

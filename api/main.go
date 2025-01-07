package main

import (
	"context"
	"encoding/json"
	"fitness-tracker/db"
	"fitness-tracker/env"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const HEADER_JWT = "jwt_token"
const JWT_EXPIRATION_TIME = 30 * time.Minute

var config env.Config

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error: %v", r)
		}
	}()
	config = env.LoadEnvVariables()
	log.Printf("main: starting backend on the port: %v", config.Env.ApiPort())
	// starting the server that will listen forever on the port
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/auth/google", corsHandler(googleAuthCodeHandler))
	http.HandleFunc("/authenticated", corsHandler(authenticatedCallHandler))
	http.HandleFunc("/logout", corsHandler(logoutCallHandler))

	// temp methods
	http.HandleFunc("/testdb", corsHandler(dbTestHandler))
	http.HandleFunc("/seeddb", corsHandler(dbSeedHandler))

	log.Fatal(http.ListenAndServe(":"+config.Env.ApiPort(), nil))
}

func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.Env.WebBaseUrl())
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, googleTokens")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")
		w.Header().Add("Content-Type", "application/json")

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
		log.Printf("main: url=%v", request.URL.Path)
	}
}

func authenticatedCallHandler(responseWriter http.ResponseWriter, request *http.Request) {
	err := validateToken(responseWriter, request)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	a, _ := json.Marshal(map[string]int{"foo": 1, "bar": 2, "baz": 3})
	responseWriter.Write(a)
}

func logoutCallHandler(responseWriter http.ResponseWriter, request *http.Request) {
	err := validateToken(responseWriter, request)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     HEADER_JWT,
		Value:    "",
		Domain:   config.Env.WebDomain(),
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false, //NOTE: for testing purposes only
	})
	a, _ := json.Marshal(map[string]bool{"logged out": true})
	responseWriter.Write(a)
}

func dbTestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	dbConnection, err := db.InitConnection(config)
	if err != nil {
		log.Printf("main: connection to db failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
	defer func() {
		dbConnection.DB.Close()
	}()

	repo := db.AuthRepo{Connection: dbConnection}

	data, err := repo.GetSomeRandomData()
	if err != nil {
		log.Printf("main: failed to get data: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("main: failed to marshal data: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
	responseWriter.Write(jsonData)

}
func dbSeedHandler(responseWriter http.ResponseWriter, request *http.Request) {
	dbConnection, err := db.InitConnection(config)
	if err != nil {
		log.Printf("main: connection to db failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
	defer func() {
		dbConnection.DB.Close()
	}()

	repo := db.AuthRepo{Connection: dbConnection}

	err = repo.Seed()
	if err != nil {
		log.Printf("main: failed to get data: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
}

/**
* Validates token consistently. Call in all authenticated calls.
 */
func validateToken(responseWriter http.ResponseWriter, request *http.Request) error {
	jwtCookie, err := request.Cookie(HEADER_JWT)
	if err != nil {
		log.Printf("main: token is missing: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return err
	}

	cryptoPublicKey := config.SecEnv.JwtKeyPublic()

	parsed, err := jwt.Parse(jwtCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return cryptoPublicKey, nil
	})

	if !parsed.Valid {
		log.Printf("main: token verification failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return err
	}
	// why do this not work while Parse above does?
	// err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), ([]byte)(parts[2]), cryptoPublicKey)
	// if err != nil {
	// 	log.Printf("Token verification failed: %v\n", err)
	// 	http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
	// 	return
	// }

	return nil
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
		log.Printf("main: reading authCode failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	var authCode GoogleAuthCodeInput
	err = json.Unmarshal(body, &authCode)
	if err != nil {
		log.Printf("main: decoding authCode failed: %v; body: %v\n", err, body)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// got authcode for Google, request tokens and store them in db
	googleTokens, err := googleOneOffTokenExchange(authCode)
	if err != nil {
		log.Printf("main: google token exchange for session failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// generate user in db and add Gtokens
	info, err := googleGetProfileInfo(googleTokens)
	if err != nil {
		log.Printf("main: google profile read failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("main: customer's info: %v\n", *info)
	//TODO: add info to DB along with refresh token

	// generating only jwt token for now. Using Google Id as subject in token
	now := time.Now()
	serverJWTToken, err := jwtWithCustomClaims(info.Id, now)
	if err != nil {
		log.Printf("main: JWT server token generated: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// set http-only jwk token cookie
	http.SetCookie(responseWriter, &http.Cookie{
		Name:     HEADER_JWT,
		Value:    serverJWTToken,
		Domain:   config.Env.WebDomain(),
		Path:     "/",
		Expires:  now.Add(JWT_EXPIRATION_TIME),
		HttpOnly: true,
		Secure:   false, //NOTE: for testing purposes only
	})
	profile, err := json.Marshal(info)
	if err != nil {
		log.Printf("main: failed to marshal profile: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
	responseWriter.Write(profile)
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

// Create a token
// Ref: https://github.com/golang-jwt/jwt/blob/main/example_test.go
func jwtWithCustomClaims(customerId string, now time.Time) (string, error) {
	cryptoKey := config.SecEnv.JwtKeyPrivate()

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
			ExpiresAt: jwt.NewNumericDate(now.Add(JWT_EXPIRATION_TIME)),
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

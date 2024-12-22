package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {

	// for now a simple load of credentials
	log.Println("Loading env variables")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!")
	}

	if os.Getenv("CLIENT_ID") == "" || os.Getenv("CLIENT_SECRET") == "" || os.Getenv("CLIENT_CALLBACK_URL") == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	fmt.Printf("Starting backend on the port: %v", httpPort)
	// starting the server that will listen forever on the port
	http.HandleFunc("/", rootHandler)
	http.Handle("/api/auth/google", corsHandler(googleAuthHandler))

	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}

func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("WEB_HOST"))
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
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

type GoogleAuthCode struct {
	Code string `json:"code"`
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func googleAuthHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var googleOauthConfig = &oauth2.Config{
		// port is of the ReactJS app
		RedirectURL:  fmt.Sprintf("http://localhost:%v", 3000),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// Declare a new Person struct.
	var authCode GoogleAuthCode

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("Reading authCode failed: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	error := json.Unmarshal(body, &authCode)
	if error != nil {
		log.Printf("Decoding authCode failed: %v; body: %v\n", error, body)
		http.Error(responseWriter, error.Error(), http.StatusBadRequest)
		return
	}

	// Use code to get token and get user info from Google.
	token, err := googleOauthConfig.Exchange(context.Background(), authCode.Code)
	if err != nil {
		log.Printf("Exchange of authcode (%v) failed: %v\n", authCode.Code, err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		log.Println("Getting profile data failed")
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Reading profile data failed")
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	responseWriter.Header().Add("Content-Type", "application/json")
	// just sending bytes
	log.Printf("Got content:\n %s\n", contents)
	fmt.Fprintf(responseWriter, "%s", contents)
}

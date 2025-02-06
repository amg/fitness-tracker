package main

import (
	"encoding/json"
	"fitness-tracker/controllers"
	"fitness-tracker/db"
	"fitness-tracker/env"
	repository "fitness-tracker/repository"
	"fitness-tracker/utils"
	"fmt"
	"io"
	"log"
	"net/http"
)

// parsed config from ENV variables and secret storage
var config env.Config

// connected DB to be used across handlers
var database *db.DB

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error: %v", r)
		}
	}()
	config = env.LoadEnvVariables()
	log.Printf("main: starting backend on the port: %v", config.Env.ApiPort())

	// initialise db connection to use across the handlers
	db, err := db.InitDB(config)
	if err != nil {
		log.Panicf("main: connection to db failed: %v\n", err)
	}
	database = db

	defer func() {
		log.Println("main: pool closed")
		database.Pool.Close()
	}()

	// starting the server that will listen forever on the port
	http.HandleFunc("/", rootHandler)

	// google auth
	http.HandleFunc("/auth/google", corsHandler(googleAuthCodeHandler))

	// session managment
	http.HandleFunc("/auth/refresh", corsHandler(authRefreshCallHandler))
	http.HandleFunc("/auth/profile", corsHandler(authGetProfileCodeHandler))
	http.HandleFunc("/auth/logout", corsHandler(authLogoutCallHandler))

	// data handlers
	http.HandleFunc("/authenticated", corsHandler(authenticatedCallHandler))

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
		responseWriter.Write([]byte(fmt.Sprintf("User Agent: %v", utils.FingerprintRequest(request))))
		responseWriter.Write([]byte(`(debug)Available endpoints:
            \t/api/auth/google
        `))
		log.Printf("main: url=%v", request.URL.Path)
	}
}

func authenticatedCallHandler(responseWriter http.ResponseWriter, request *http.Request) {
	userId, err := utils.ValidateSession(config, request)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	a, _ := json.Marshal(map[string]any{"userId": userId.String()})
	responseWriter.Write(a)
}

func authLogoutCallHandler(responseWriter http.ResponseWriter, request *http.Request) {
	userId, err := utils.ValidateSession(config, request)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	authRepo := repository.AuthRepo{DB: database}
	// TODO: wrap with controller to refrain from calling repo directly from handlers
	authRepo.DeleteRefreshTokenByUserId(*userId, utils.FingerprintRequest(request))
	// reset cookies
	utils.ClearTokenCookies(config, responseWriter, request)

	a, _ := json.Marshal(map[string]bool{"logged out": true})
	responseWriter.Write(a)
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

	var authCode controllers.GoogleAuthCodeInput
	err = json.Unmarshal(body, &authCode)
	if err != nil {
		log.Printf("main: decoding authCode failed: %v; body: %v\n", err, body)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	profile, err := controllers.AuthenticateWithGoogle(config, authCode)
	if err != nil {
		log.Printf("main: Google auth failed: %v", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// create or merge customer
	authRepo := repository.AuthRepo{DB: database}
	customerRecord, tokens, err := controllers.CreateOrMergeCustomer(
		config, &authRepo, utils.FingerprintRequest(request),
		profile.Email, profile.GivenName, profile.FamilyName, profile.Picture)
	if err != nil {
		log.Printf("main: failed to create new customer and issue tokens: %v", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// set http-only jwk token cookie
	utils.SetTokenCookies(config, tokens.Session, tokens.Refresh, responseWriter, request)

	customerJson, err := json.Marshal(customerRecord)
	if err != nil {
		log.Printf("main: failed to marshal customer: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.Write(customerJson)
}

func authGetProfileCodeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	userId, err := utils.ValidateSession(config, request)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	authRepo := repository.AuthRepo{DB: database}
	// TODO: wrap with controller to refrain from calling repo directly from handlers
	customerRecord, err := authRepo.GetCustomerInfo(*userId)
	if err != nil {
		log.Printf("main: failed to find customer: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusForbidden)
	}

	customerJson, err := json.Marshal(customerRecord)
	if err != nil {
		log.Printf("main: failed to marshal customer: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
	responseWriter.Write(customerJson)
}

func authRefreshCallHandler(responseWriter http.ResponseWriter, request *http.Request) {
	_, refreshToken, err := utils.ValidateRefreshToken(config, request)
	if err != nil {
		log.Printf("main: invalid token: %v", err)
		http.Error(responseWriter, err.Error(), http.StatusForbidden)
		return
	}

	authRepo := repository.AuthRepo{DB: database}
	session, refresh, err := controllers.ReIssueTokens(config, &authRepo, utils.FingerprintRequest(request), refreshToken)
	if err != nil {
		log.Printf("main: failed to find refresh token: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusForbidden)
		return
	}
	// set http-only jwk token cookie
	utils.SetTokenCookies(config, session, refresh, responseWriter, request)
}

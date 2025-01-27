package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fitness-tracker/env"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

/**
* Validates token consistently. Call in all authenticated calls.
 */
func ValidateSession(config env.Config, responseWriter http.ResponseWriter, request *http.Request) (userId *uuid.UUID, err error) {
	jwtCookie, err := request.Cookie(KEY_SESSION_TOKEN)
	if err != nil {
		log.Printf("main: token is missing: %v\n", err)
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return nil, err
	}

	userId, err = ValidateToken(config, jwtCookie.Value)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return nil, err
	}
	return
}

// Validates refresh token (JWT) and returns userId and token value
func ValidateRefreshToken(config env.Config, responseWriter http.ResponseWriter, request *http.Request) (userId *uuid.UUID, refreshToken string, err error) {
	refreshTokenCookie, err := request.Cookie(KEY_REFRESH_TOKEN)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return nil, "", err
	}

	userId, err = ValidateToken(config, refreshTokenCookie.Value)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
		return nil, "", err
	}

	return userId, refreshTokenCookie.Value, nil
}

// Create a unique user's device fingerprint.
// Refresh tokens are generated uniquely on per fingerprint basis
func FingerprintRequest(request *http.Request) string {
	hasher := sha256.New()
	hasher.Write(([]byte)(request.UserAgent()))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func SetTokenCookies(config env.Config, session string, refresh string, responseWriter http.ResponseWriter, request *http.Request) {
	var secure bool
	switch config.Env.(type) {
	case env.DevEnv:
		secure = false
	case env.StagingEnv:
		secure = true
	default:
		panic("db: unsupported env")
	}

	now := time.Now()
	http.SetCookie(responseWriter, &http.Cookie{
		Name:     KEY_SESSION_TOKEN,
		Value:    session,
		Domain:   config.Env.WebDomain(),
		Path:     "/",
		Expires:  now.Add(SESSION_EXPIRATION_TIME),
		HttpOnly: true,
		Secure:   secure,
	})
	http.SetCookie(responseWriter, &http.Cookie{
		Name:     KEY_REFRESH_TOKEN,
		Value:    refresh,
		Domain:   config.Env.WebDomain(),
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   secure,
	})
}

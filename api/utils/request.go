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
func ValidateSession(config env.Config, request *http.Request) (userId *uuid.UUID, err error) {
	jwtCookie, err := request.Cookie(KEY_SESSION_TOKEN)
	if err != nil {
		log.Printf("main: token is missing: %v\n", err)
		return nil, err
	}

	userId, _, err = ValidateToken(config, jwtCookie.Value)
	if err != nil {
		return nil, err
	}
	return
}

// Validates refresh token (JWT) and returns userId and token value
func ValidateRefreshToken(config env.Config, request *http.Request) (userId *uuid.UUID, refreshToken string, err error) {
	refreshTokenCookie, err := request.Cookie(KEY_REFRESH_TOKEN)
	if err != nil {
		return nil, "", err
	}

	userId, _, err = ValidateToken(config, refreshTokenCookie.Value)
	if err != nil {
		return nil, "", err
	}

	return userId, refreshTokenCookie.Value, nil
}

// Create user's device fingerprint.
// Refresh tokens are generated securely and uniquely and fingerprint is stored along side for future.
// Fingerprint hashing is not really that secure as format for UserAgent is known and it's easy to enumerate all of the possible values
// even if it were salted, it will not really slow down reverse engineering but data inside is not really a secret.
//
// For now this will suffice
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

func ClearTokenCookies(config env.Config, responseWriter http.ResponseWriter, request *http.Request) {
	var secure bool
	switch config.Env.(type) {
	case env.DevEnv:
		secure = false
	case env.StagingEnv:
		secure = true
	default:
		panic("db: unsupported env")
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     KEY_SESSION_TOKEN,
		Value:    "",
		Domain:   config.Env.WebDomain(),
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   secure,
	})
	http.SetCookie(responseWriter, &http.Cookie{
		Name:     KEY_REFRESH_TOKEN,
		Value:    "",
		Domain:   config.Env.WebDomain(),
		Path:     "/auth/refresh",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   secure,
	})
}

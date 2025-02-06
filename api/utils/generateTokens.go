package utils

import (
	"fitness-tracker/env"
	"log"
	"time"

	"github.com/google/uuid"
)

// Generate new pair of session and refresh tokens for the customer
func GenerateSessionToken(config env.Config, customerId string) (sessionToken string, err error) {
	// generating only jwt token for now. Using Google Id as subject in token
	now := time.Now()
	sessionToken, _, err = JwtWithCustomClaims(config, customerId, now, now.Add(SESSION_EXPIRATION_TIME))
	if err != nil {
		log.Printf("main: JWT server token generation failed: %v\n", err)
		return "", err
	}

	return
}

// Generate new pair of session and refresh tokens for the customer
func GenerateRefreshToken(config env.Config, customerId string) (refreshToken string, jti *uuid.UUID, expiresAt time.Time, err error) {
	// generating only jwt token for now. Using Google Id as subject in token
	now := time.Now()
	expiration := now.Add(REFRESH_TOKEN_EXPIRATION_TIME)
	refreshToken, jti, err = JwtWithCustomClaims(config, customerId, now, expiration)
	if err != nil {
		log.Printf("main: refresh token generation failed: %v\n", err)
		return "", nil, now, err
	}

	return refreshToken, jti, expiration, nil
}

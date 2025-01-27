package utils

import (
	"fitness-tracker/env"
	"log"
	"time"
)

// Generate new pair of session and refresh tokens for the customer
func GenerateTokens(config env.Config, customerId string) (sessionToken string, refreshToken string, err error) {
	// generating only jwt token for now. Using Google Id as subject in token
	now := time.Now()
	sessionToken, err = JwtWithCustomClaims(config, customerId, now, now.Add(SESSION_EXPIRATION_TIME))
	if err != nil {
		log.Printf("main: JWT server token generation failed: %v\n", err)
		return "", "", err
	}

	refreshToken, err = JwtWithCustomClaims(config, customerId, now, now.Add(REFRESH_TOKEN_EXPIRATION_TIME))
	if err != nil {
		log.Printf("main: refresh token generation failed: %v\n", err)
		return "", "", err
	}

	return
}

package utils

import (
	"fitness-tracker/env"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Create a token
// Ref: https://github.com/golang-jwt/jwt/blob/main/example_test.go
func JwtWithCustomClaims(config env.Config, customerId string, now time.Time, expiresAt time.Time) (string, error) {
	cryptoKey := config.SecEnv.JwtKeyPrivate()

	type MyCustomClaims struct {
		jwt.RegisteredClaims
	}

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
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

/**
* Validates token consistently. Call in all authenticated calls.
 */
func ValidateToken(config env.Config, token string) (userId *uuid.UUID, err error) {
	cryptoPublicKey := config.SecEnv.JwtKeyPublic()

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", token.Header["alg"])
		}

		return cryptoPublicKey, nil
	})

	if !parsed.Valid {
		log.Printf("jwt: token verification failed: %v\n", err)
		return nil, err
	}

	userIdString, err := parsed.Claims.GetSubject()
	if err != nil {
		log.Printf("jwt: couldn't get subject: %v", err)
		return nil, err
	}

	user, err := uuid.Parse(userIdString)
	if err != nil {
		log.Printf("jwt: couldn't parse into uuid: %v", err)
		return nil, err
	}
	userId = &user

	return
}

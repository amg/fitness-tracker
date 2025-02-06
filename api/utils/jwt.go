package utils

import (
	"fitness-tracker/env"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MyCustomClaims struct {
	jwt.RegisteredClaims
}

// Create a token
// Ref: https://github.com/golang-jwt/jwt/blob/main/example_test.go
func JwtWithCustomClaims(config env.Config, customerId string, now time.Time, expiresAt time.Time) (signedJwt string, id *uuid.UUID, err error) {
	cryptoKey := config.SecEnv.JwtKeyPrivate()

	uuid, err := uuid.NewV7()
	if err != nil {
		return "", nil, fmt.Errorf("tokens: failed to generate uuidv7; %v", err)
	}
	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "FitnessTracker",
			Subject:   customerId,
			// generate pseudo random almost unique UUID
			// for session tokens they don't have to be unique, customerId will be
			// for refresh tokens they will need to be unique as this value is stored in DB
			//  to allow blacklisting. This is controlled DB UNIQUE and backoff logic to retry in
			//  extremely rare event of collision
			ID:       uuid.String(),
			Audience: []string{"FitnessTrackerAPI"},
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)

	signedJwt, err = token.SignedString(cryptoKey)
	if err != nil {
		return "", nil, fmt.Errorf("tokens: failed to sign token; %v", err)
	}

	return signedJwt, &uuid, nil
}

// Validates token consistently. Call in all authenticated calls.
// `userId` is unique customer internal identifier
// `id` is unique token identifier
// `err` error
func ValidateToken(config env.Config, token string) (userId *uuid.UUID, id *uuid.UUID, err error) {
	cryptoPublicKey := config.SecEnv.JwtKeyPublic()

	parsed, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", token.Header["alg"])
		}

		return cryptoPublicKey, nil
	})

	if !parsed.Valid {
		return nil, nil, fmt.Errorf("jwt: token verification failed: %v", err)
	}

	allClaims := parsed.Claims.(*MyCustomClaims)
	userIdString, err := allClaims.GetSubject()
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: couldn't get subject: %v", err)
	}

	if allClaims.ID == "" {
		return nil, nil, fmt.Errorf("jwt: couldn't get jti: %v", err)
	}

	user, err := uuid.Parse(userIdString)
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: couldn't parse 'userId' into uuid: %v", err)
	}
	jti, err := uuid.Parse(allClaims.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: couldn't parse 'jti' into uuid: %v", err)
	}

	return &user, &jti, nil
}

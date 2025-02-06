package controllers

import (
	"context"
	"fitness-tracker/env"
	"fitness-tracker/models"
	"fitness-tracker/repository"
	"fitness-tracker/utils"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func CreateOrMergeCustomer(
	config env.Config,
	authRepo *repository.AuthRepo,
	fingerprint string,
	email string,
	firstName string,
	lastName string,
	picture string,
) (user *models.UserInfo, tokens *models.JWTokens, err error) {
	customerRecord, err := authRepo.CreateOrMergeCustomer(email, firstName, lastName, picture)
	if err != nil {
		return nil, nil, fmt.Errorf("authController: creating/finding customer failed: %v", err)
	}

	session, refresh, err := issueTokens(config, authRepo, customerRecord.ID, fingerprint)
	if err != nil {
		return nil, nil, fmt.Errorf("authController: failed to marshal customer: %v", err)
	}

	return customerRecord, &models.JWTokens{Session: session, Refresh: refresh}, nil
}

// Find refresh token in the db and reissue session and refresh
// Stores new refresh token in the db (deletes previous one)
func ReIssueTokens(config env.Config, authRepo *repository.AuthRepo, fingerprint string, currentRefreshToken string) (sessionToken string, refreshToken string, err error) {
	ctx := context.Background()

	// validate and parse refresh jwt to get user id
	_, jti, err := utils.ValidateToken(config, currentRefreshToken)
	if err != nil {
		// validation failed
		return "", "", fmt.Errorf("authController: the refresh token is invalid; %v", err)
	}

	foundToken, err := authRepo.GetRefreshToken(*jti)
	if err != nil {
		return "", "", fmt.Errorf("authController: refresh token (jti) is not found; %v", err)
	}

	err = authRepo.DB.Queries.DeleteRefreshToken(ctx, foundToken.ID)
	if err != nil {
		// this should never happen as we have presumably found token above
		return "", "", fmt.Errorf("authController: couldn't delete refresh token; %v", err)
	}

	return issueTokens(config, authRepo, foundToken.UserID, fingerprint)
}

func issueTokens(config env.Config, authRepo *repository.AuthRepo, userId uuid.UUID, fingerprint string) (sessionToken string, refreshToken string, err error) {
	// issue tokens
	session, err := utils.GenerateSessionToken(config, userId.String())
	if err != nil {
		return "", "", fmt.Errorf("authController: failed to generate session token: %v", err)
	}

	// try up to 3 times to generate a unique token
	// UUID is plenty unique but DB requires strict unique, so just in case loop up to 3 times
	var refresh string
	for i := 1; i <= 3; i++ {
		var expiresAt time.Time
		var jti *uuid.UUID
		refresh, jti, expiresAt, err = utils.GenerateRefreshToken(config, userId.String())
		if err != nil {
			return "", "", fmt.Errorf("authController: failed to generate refresh token: %v", err)
		}

		_, err = authRepo.UpsertRefreshToken(userId, fingerprint, *jti, expiresAt)
		if err != nil {
			log.Printf("authController: token collision detected; Attempt %v out of %v", i, 3)
			refresh = ""
		} else {
			break
		}
	}
	// if still empty and didn't get upserted
	if refresh == "" {
		return "", "", fmt.Errorf("authController: failed to generate unique token; %v", err)
	}

	return session, refresh, nil
}

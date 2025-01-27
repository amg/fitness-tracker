package controllers

import (
	"context"
	"fitness-tracker/env"
	"fitness-tracker/repository"
	"fitness-tracker/utils"
	"log"
)

// Find refresh token in the db and reissue session and refresh
// Stores new refresh token in the db (deletes previous one)
func ReIssueTokens(config env.Config, authRepo *repository.AuthRepo, currentRefreshToken string) (sessionToken string, refreshToken string, err error) {
	ctx := context.Background()
	foundToken, err := authRepo.GetRefreshToken(currentRefreshToken)
	if err != nil {
		log.Println("authController: couldn't find the refresh token")
		return "", "", err
	}
	err = authRepo.DB.Queries.DeleteRefreshToken(ctx, foundToken.ID)
	if err != nil {
		log.Println("authController: couldn't delete token")
		return "", "", err
	}

	session, refresh, err := utils.GenerateTokens(config, foundToken.UserID.String())
	if err != nil {
		log.Printf("authController: failed to generate tokens: %v", err)
		return "", "", err
	}

	// use old fingerprint
	_, err = authRepo.UpsertRefreshToken(foundToken.UserID, foundToken.Fingerprint, refresh)
	if err != nil {
		log.Printf("authController: failed to generate refresh token: %v", err)
		return "", "", err
	}

	return session, refresh, nil
}

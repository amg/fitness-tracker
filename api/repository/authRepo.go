package repository

import (
	"context"

	"fitness-tracker/db"
	gendb "fitness-tracker/db/generated"
	"fitness-tracker/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type AuthRepo struct {
	DB *db.DB
}

type RandomData struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func (authRepo *AuthRepo) CreateOrMergeCustomer(email string, firstName string, lastName string, picture string) (user *models.UserInfo, err error) {
	ctx := context.Background()

	userInfo, err := authRepo.DB.Queries.GetUserInfoByEmail(ctx, email)
	// no error means user is found with this email
	if err == nil {
		return &models.UserInfo{
			ID:         userInfo.ID,
			Email:      userInfo.Email,
			FirstName:  userInfo.FirstName,
			LastName:   userInfo.LastName,
			PictureUrl: *userInfo.PictureUrl,
		}, nil
	}

	// no user found, creating new one
	userInfo, err = authRepo.DB.Queries.CreateUserInfo(ctx, gendb.CreateUserInfoParams{
		ID:         uuid.New(),
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		PictureUrl: &picture,
	})
	if err != nil {
		return nil, err
	}
	var pic string
	if userInfo.PictureUrl == nil {
		pic = ""
	} else {
		pic = *userInfo.PictureUrl
	}
	return &models.UserInfo{
		ID:         userInfo.ID,
		Email:      userInfo.Email,
		FirstName:  userInfo.FirstName,
		LastName:   userInfo.LastName,
		PictureUrl: pic,
	}, nil
}

func (authRepo *AuthRepo) UpsertRefreshToken(userId uuid.UUID, fingerprint string, refreshToken string) (token string, err error) {
	updatedToken, err := authRepo.DB.Queries.CreateRefreshToken(context.Background(), gendb.CreateRefreshTokenParams{
		ID:          refreshToken,
		UserID:      userId,
		Fingerprint: fingerprint,
	})
	if err != nil {
		return "", err
	}
	return updatedToken.ID, nil
}

func (authRepo *AuthRepo) DeleteRefreshTokenByUserId(userId uuid.UUID, fingerprint string) error {
	return authRepo.DB.Queries.DeleteRefreshTokenByUserAndFingerprint(context.Background(), gendb.DeleteRefreshTokenByUserAndFingerprintParams{
		UserID:      userId,
		Fingerprint: fingerprint,
	})
}

func (authRepo *AuthRepo) GetCustomerInfo(userId uuid.UUID) (userInfo *models.UserInfo, err error) {
	ctx := context.Background()
	info, err := authRepo.DB.Queries.GetUserInfo(ctx, userId)
	// no error means user is found with this email
	if err != nil {
		return nil, err
	}
	return &models.UserInfo{
		ID:         info.ID,
		Email:      info.Email,
		FirstName:  info.FirstName,
		LastName:   info.LastName,
		PictureUrl: *info.PictureUrl,
	}, nil
}

func (authRepo *AuthRepo) GetRefreshToken(currentToken string) (refreshToken *models.RefreshToken, err error) {
	ctx := context.Background()
	foundToken, err := authRepo.DB.Queries.GetRefreshToken(ctx, currentToken)
	if err != nil {
		return nil, err
	}
	return &models.RefreshToken{
		ID:          foundToken.ID,
		UserID:      foundToken.UserID,
		Fingerprint: foundToken.Fingerprint,
		CreatedAt:   foundToken.CreatedAt.Time,
	}, nil
}

func (authRepo *AuthRepo) DeleteRefreshToken(refreshToken string) error {
	ctx := context.Background()
	return authRepo.DB.Queries.DeleteRefreshToken(ctx, refreshToken)
}

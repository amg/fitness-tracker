package utils

import (
	"crypto/rsa"
	"fitness-tracker/env"
	"os"
	"testing"
	"time"
)

const fiftyYears = 50 * 365 * 24 * time.Hour

type MockNonSecEnv struct{}

func (m *MockNonSecEnv) GoogleClientId() string {
	return ""
}
func (m *MockNonSecEnv) GoogleClientCallbackUrl() string {
	return ""
}
func (m *MockNonSecEnv) ApiPort() string {
	return ""
}
func (m *MockNonSecEnv) WebDomain() string {
	return ""
}
func (m *MockNonSecEnv) WebBaseUrl() string {
	return ""
}
func (m *MockNonSecEnv) PostgresUrl() string {
	return ""
}
func (m *MockNonSecEnv) GoogleClientSecret() string {
	return ""
}
func (m *MockNonSecEnv) JwtKeyPrivate() *rsa.PrivateKey {
	return env.LoadRSAPrivateKeyFromDisk("../testsdata/test-jwt-rsa256-private.pem")
}
func (m *MockNonSecEnv) JwtKeyPublic() *rsa.PublicKey {
	return env.LoadRSAPublicKeyFromDisk("../testsdata/test-jwt-rsa256-public.pem")
}
func (m *MockNonSecEnv) PostgresDbName() string {
	return ""
}
func (m *MockNonSecEnv) PostgresUser() string {
	return ""
}
func (m *MockNonSecEnv) PostgresPassword() string {
	return ""
}

func Test_JwtWithCustomerClaims_Generated_Valid(t *testing.T) {
	var testTime time.Time
	testTime, _ = time.Parse(time.RFC3339, "2024-01-02T15:04:05Z")
	mockEnv := MockNonSecEnv{}
	mockConfig := env.Config{Env: &mockEnv, SecEnv: &mockEnv}

	userIdString := "3626460a-4c5d-4779-8e04-5180851c2f5c"
	jwt, _, _ := JwtWithCustomClaims(mockConfig, userIdString, testTime, testTime.Add(fiftyYears))

	userId, id, err := ValidateToken(mockConfig, jwt)
	if err != nil {
		t.Errorf("failed to validate token: '%v'", jwt)
	}
	if userId.String() != userIdString {
		t.Errorf("failed to encode user id into the token: '%v'", jwt)
	}
	if id == nil {
		t.Errorf("failed to generate/encode id into the token: '%v'", jwt)
	}
}

func Test_ValidateToken_Valid(t *testing.T) {
	mockEnv := MockNonSecEnv{}
	mockConfig := env.Config{Env: &mockEnv, SecEnv: &mockEnv}

	validToken, err := os.ReadFile("../testsdata/validToken.txt")
	if err != nil {
		t.Error(err)
	}
	_, _, err = ValidateToken(mockConfig, string(validToken))
	if err != nil {
		t.Error(err)
	}
}

func Test_ValidateToken_Invalid(t *testing.T) {
	mockEnv := MockNonSecEnv{}
	mockConfig := env.Config{Env: &mockEnv, SecEnv: &mockEnv}

	validToken, err := os.ReadFile("../testsdata/validToken.txt")
	if err != nil {
		t.Error(err)
	}
	_, _, err = ValidateToken(mockConfig, string(validToken)+"corrupt signature")
	if err == nil {
		t.Error("Signature is not correctly verified")
	}
}

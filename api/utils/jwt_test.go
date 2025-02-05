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

	jwt, _ := JwtWithCustomClaims(mockConfig, "3626460a-4c5d-4779-8e04-5180851c2f5c", testTime, testTime.Add(fiftyYears))
	expected := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJpc3MiOiJGaXRuZXNzVHJhY2tlciIsInN1YiI6IjM2MjY0NjBhLTRjNWQtNDc3OS04ZTA0LTUxODA4NTFjMmY1YyIsImF1ZCI6WyJGaXRuZXNzVHJhY2tlckFQSSJdLCJleHAiOjMyODEwMDc4NDUsIm5iZiI6MTcwNDIwNzg0NSwiaWF0IjoxNzA0MjA3ODQ1fQ.CEcyPR57YwzuPFG1O9_Ls716b3UVGDhVeMgckKWsxBj2m5zjqO6YqKD0xedn1h0n8qBr_81rXo98wZnnwr699fFZcfkVQtCGwXUws_AYr-JoBT6X3uvQBez-9TpkmIBGlkjtKMo83WaLbFTxTCXlWIwQ4cwmemnz6GV9xK7rE-81eqSO1zdlldY2Wqa3hDEg8MnlmOx7NViXWmxE0kpg18W-TUydl9LJ_jSfcUCFUmCCeblcNvHFy3jVFpgWUwNf_bsOIsiDWvIPTmm-bZE_UqItC7kU4ZXhecvJBkeJFoSNjZV22WdYSMNckbPT-bVg-gOqyCMe4a5GIEjLuQUnqA"
	if expected != jwt {
		t.Error(jwt)
	}
}

func Test_ValidateToken_Valid(t *testing.T) {
	mockEnv := MockNonSecEnv{}
	mockConfig := env.Config{Env: &mockEnv, SecEnv: &mockEnv}

	validToken, err := os.ReadFile("../testsdata/validToken.txt")
	if err != nil {
		t.Error(err)
	}
	_, err = ValidateToken(mockConfig, string(validToken))
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
	_, err = ValidateToken(mockConfig, string(validToken)+"corrupt signature")
	if err == nil {
		t.Error("Signature is not correctly verified")
	}
}

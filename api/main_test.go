package main

import (
	"crypto/rsa"
	"fitness-tracker/env"
	"fmt"
	"testing"
	"time"
)

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
	// TODO: shouldn't depend on another service's test files but for now it will do
	return env.LoadRSAPrivateKeyFromDisk("../node-api/__tests__/src/services/test-jwt-rsa256-private.pem")
}
func (m *MockNonSecEnv) JwtKeyPublic() *rsa.PublicKey {
	return nil
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

func TestJWTGenerationAndSigning(t *testing.T) {
	var testTime time.Time
	testTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	fmt.Println(testTime)
	mockEnv := MockNonSecEnv{}
	mockConfig := env.Config{Env: &mockEnv, SecEnv: &mockEnv}

	jwt, _ := jwtWithCustomClaims(mockConfig, "../.secrets/jwtRSA256-private.pem", testTime)
	expected := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJpc3MiOiJGaXRuZXNzVHJhY2tlciIsInN1YiI6Ii4uLy5zZWNyZXRzL2p3dFJTQTI1Ni1wcml2YXRlLnBlbSIsImF1ZCI6WyJGaXRuZXNzVHJhY2tlckFQSSJdLCJleHAiOjExMzYyMTYwNDUsIm5iZiI6MTEzNjIxNDI0NSwiaWF0IjoxMTM2MjE0MjQ1fQ.eZmuHICGSlz4KJ-vQGI74Qy7oZwkLKXNDhmImtzzHxrn-IK8y77A-nb37sFyZ591WfSRGvJrlcUSXE_Ijrc9CN688lu5fWYKMs_68e0S0fEQcU59BrLYR9lvvyANzEixo0RgtfNcdht84-VHtZVQJAEHJiHSacLTWw_rC5hgNURliWk6VpUYagXXYUe8-9KtQrIFFspAj1E7BM60KVzUTFHuAyQ2HHSoAtm2KRL6AuRlrqQJqAQNqkbs4uOQCGw7K2Tzs9Vuab6w7vFnXdzT9SD3TcQP5XF3SFPfM2p8pafrOvg3OYpyBzTsrSLG1DWGYtuOCFF4Ecs1nHgeQyqTcQ"
	if expected != jwt {
		t.Error(jwt)
	}
}

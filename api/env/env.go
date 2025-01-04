package env

import (
	"context"
	"crypto/rsa"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"sync"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

/**
* Use for non sensitive env variables
 */
type NonSecureEnv interface {
	GoogleClientId() string
	GoogleClientCallbackUrl() string
	ApiPort() string
	WebDomain() string
	WebBaseUrl() string
}

/**
* Use for secure variables that should not be exposed in production
 */
type SecureEnv interface {
	GoogleClientSecret() string
	JwtKeyPrivate() *rsa.PrivateKey
	JwtKeyPublic() *rsa.PublicKey
}

type Config struct {
	Env    NonSecureEnv
	SecEnv SecureEnv
}

func LoadEnvVariables() *Config {
	log.Println("Loading env variables")
	return sync.OnceValue(doLoadEnv)()
}

func doLoadEnv() *Config {
	log.Println("Do Loading env variables")

	googleProjectId := loadFromEnv("GOOGLE_PROJECT_ID")
	googleClientId := loadFromEnv("GOOGLE_CLIENT_ID")
	googleClientCallbackUrl := loadFromEnv("GOOGLE_CLIENT_CALLBACK_URL")
	apiPort := loadFromEnv("API_PORT")
	webDomain := loadFromEnv("COOKIE_DOMAIN")
	webBaseUrl := loadFromEnv("WEB_BASE_URL")

	log.Printf("Loaded basic env variables; env:%v", os.Getenv("ENV"))
	switch os.Getenv("ENV") {
	case "dev":
		devEnv := DevEnv{
			Env{
				googleClientId:          googleClientId,
				googleClientCallbackUrl: googleClientCallbackUrl,
				apiPort:                 apiPort,
				webDomain:               webDomain,
				webBaseUrl:              webBaseUrl,
				googleClientSecret:      loadFromEnv("GOOGLE_CLIENT_SECRET"),
				jwtKeyPrivate:           LoadRSAPrivateKeyFromDisk(loadFromEnv("FILE_KEY_PRIVATE")),
				jwtKeyPublic:            LoadRSAPublicKeyFromDisk(loadFromEnv("FILE_KEY_PUBLIC")),
			},
		}
		log.Printf("Created config. Web url: %v", devEnv.webBaseUrl)
		return &Config{Env: &devEnv, SecEnv: &devEnv}
	case "staging":
		ctx := context.Background()
		client, err := secretmanager.NewClient(ctx)
		if err != nil {
			log.Panicf("Failed to setup client: %v", err)
		}
		defer client.Close()

		googleOAuthSecret := secretByKey("GOOGLE_OAUTH_CLIENT_SECRET", client, googleProjectId)
		jwtKeyPrivateString := secretByKey("JWT_KEY_PRIVATE", client, googleProjectId)
		jwtKeyPrivate := ParseRSAPrivateKeyFromPEMString(([]byte)(jwtKeyPrivateString))

		jwtKeyPublicString := secretByKey("JWT_KEY_PUBLIC", client, googleProjectId)
		jwtKeyPublic := ParseRSAPublicKeyFromPEMString(([]byte)(jwtKeyPublicString))

		devEnv := StagingEnv{
			Env{
				googleClientId:          googleClientId,
				googleClientCallbackUrl: googleClientCallbackUrl,
				apiPort:                 apiPort,
				webDomain:               webDomain,
				webBaseUrl:              webBaseUrl,
				googleClientSecret:      googleOAuthSecret,
				jwtKeyPrivate:           jwtKeyPrivate,
				jwtKeyPublic:            jwtKeyPublic,
			},
		}
		return &Config{Env: &devEnv, SecEnv: &devEnv}
	default:
		log.Println("Failed")
		panic("Unsupported env")
	}
}

func secretByKey(key string, client *secretmanager.Client, googleProjectId string) string {
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v/versions/latest", googleProjectId, key),
	}

	// Call the API.
	result, err := client.AccessSecretVersion(context.Background(), accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		log.Fatalf("Data corruption detected retrieving JWT_KEY_PRIVATE")
	}

	return string(result.Payload.Data)
}

type Env struct {
	googleClientId          string
	googleClientCallbackUrl string
	apiPort                 string
	webDomain               string
	webBaseUrl              string
	googleClientSecret      string
	jwtKeyPrivate           *rsa.PrivateKey
	jwtKeyPublic            *rsa.PublicKey
}

type DevEnv struct {
	Env
}

type StagingEnv struct {
	Env
}

func (env Env) GoogleClientId() string {
	return env.googleClientId
}
func (env Env) GoogleClientCallbackUrl() string {
	return env.googleClientCallbackUrl
}
func (env Env) ApiPort() string {
	return env.apiPort
}
func (env Env) WebDomain() string {
	return env.webDomain
}
func (env Env) WebBaseUrl() string {
	return env.webBaseUrl
}

func (devEnv DevEnv) GoogleClientSecret() string {
	return devEnv.googleClientSecret
}

func (devEnv DevEnv) JwtKeyPrivate() *rsa.PrivateKey {
	return devEnv.jwtKeyPrivate
}
func (devEnv DevEnv) JwtKeyPublic() *rsa.PublicKey {
	return devEnv.jwtKeyPublic
}

/**
* Loads from env by key or panics
 */
func loadFromEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("%v needs to be defined", key)
	}
	return value
}

// ---------- Staging ------------
func (stagingEnv StagingEnv) GoogleClientSecret() string {
	return stagingEnv.googleClientSecret
}

func (stagingEnv StagingEnv) JwtKeyPrivate() *rsa.PrivateKey {
	return stagingEnv.jwtKeyPrivate
}
func (stagingEnv StagingEnv) JwtKeyPublic() *rsa.PublicKey {
	return stagingEnv.jwtKeyPublic
}

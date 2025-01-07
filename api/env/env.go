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
	PostgresUrl() string
}

/**
* Use for secure variables that should not be exposed in production
 */
type SecureEnv interface {
	GoogleClientSecret() string
	JwtKeyPrivate() *rsa.PrivateKey
	JwtKeyPublic() *rsa.PublicKey
	PostgresDbName() string
	PostgresUser() string
	PostgresPassword() string
}

type Config struct {
	Env    NonSecureEnv
	SecEnv SecureEnv
}

/**
* Can return either pointer or struct.
* A bit unclear if it matters much for small object like this.
* Seems simpler not to deal with pointers though.
 */
func LoadEnvVariables() Config {
	log.Println("env: loading env variables")
	return sync.OnceValue(doLoadEnv)()
}

func doLoadEnv() Config {
	log.Println("env: once Loading env variables")

	googleProjectId := loadFromEnv("GOOGLE_PROJECT_ID")
	googleClientId := loadFromEnv("GOOGLE_CLIENT_ID")
	googleClientCallbackUrl := loadFromEnv("GOOGLE_CLIENT_CALLBACK_URL")
	apiPort := loadFromEnv("API_PORT")
	webDomain := loadFromEnv("COOKIE_DOMAIN")
	webBaseUrl := loadFromEnv("WEB_BASE_URL")

	log.Printf("env: loaded basic env: %v", os.Getenv("ENV"))
	switch os.Getenv("ENV") {
	case "dev":
		postgresUrl := loadFromEnv("POSTGRES_URL")
		devEnv := DevEnv{
			Env{
				googleClientId:          googleClientId,
				googleClientCallbackUrl: googleClientCallbackUrl,
				apiPort:                 apiPort,
				webDomain:               webDomain,
				webBaseUrl:              webBaseUrl,
				postgresUrl:             postgresUrl,
				googleClientSecret:      loadFromEnv("GOOGLE_CLIENT_SECRET"),
				jwtKeyPrivate:           LoadRSAPrivateKeyFromDisk(loadFromEnv("FILE_KEY_PRIVATE")),
				jwtKeyPublic:            LoadRSAPublicKeyFromDisk(loadFromEnv("FILE_KEY_PUBLIC")),
				postgresDbName:          loadFromEnv("POSTGRES_DBNAME"),
				postgresUser:            loadFromEnv("POSTGRES_USER"),
				postgresPassword:        loadFromEnv("POSTGRES_PASSWORD"),
			},
		}
		log.Printf("env: created config. Web url: %v", devEnv.webBaseUrl)
		// could use pointers here but then less obvious when comparing .type using switch
		return Config{Env: devEnv, SecEnv: devEnv}
	case "staging":
		// for staging connection name is provided by cloud run instance
		postgresUrl := loadFromEnv("DB_INSTANCE_CONNECTION_NAME")

		ctx := context.Background()
		client, err := secretmanager.NewClient(ctx)
		if err != nil {
			log.Panicf("env: failed to setup client: %v", err)
		}
		defer client.Close()

		googleOAuthSecret := secretByKey("GOOGLE_OAUTH_CLIENT_SECRET", client, googleProjectId)
		jwtKeyPrivateString := secretByKey("JWT_KEY_PRIVATE", client, googleProjectId)
		jwtKeyPrivate := ParseRSAPrivateKeyFromPEMString(([]byte)(jwtKeyPrivateString))

		jwtKeyPublicString := secretByKey("JWT_KEY_PUBLIC", client, googleProjectId)
		jwtKeyPublic := ParseRSAPublicKeyFromPEMString(([]byte)(jwtKeyPublicString))

		postgresDbName := secretByKey("POSTGRES_DBNAME", client, googleProjectId)
		postgresUser := secretByKey("POSTGRES_USER", client, googleProjectId)
		postgresPassword := secretByKey("POSTGRES_PASSWORD", client, googleProjectId)

		devEnv := StagingEnv{
			Env{
				googleClientId:          googleClientId,
				googleClientCallbackUrl: googleClientCallbackUrl,
				apiPort:                 apiPort,
				webDomain:               webDomain,
				webBaseUrl:              webBaseUrl,
				postgresUrl:             postgresUrl,
				googleClientSecret:      googleOAuthSecret,
				jwtKeyPrivate:           jwtKeyPrivate,
				jwtKeyPublic:            jwtKeyPublic,
				postgresDbName:          postgresDbName,
				postgresUser:            postgresUser,
				postgresPassword:        postgresPassword,
			},
		}
		// could use pointers here but then less obvious when comparing .type using switch
		return Config{Env: devEnv, SecEnv: devEnv}
	default:
		panic(fmt.Sprintf("env: unsupported env: %v", os.Getenv("ENV")))
	}
}

func secretByKey(key string, client *secretmanager.Client, googleProjectId string) string {
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v/versions/latest", googleProjectId, key),
	}

	// Call the API.
	result, err := client.AccessSecretVersion(context.Background(), accessRequest)
	if err != nil {
		log.Fatalf("env: failed to access secret version: %v", err)
	}

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		log.Fatalf("env: data corruption detected retrieving JWT_KEY_PRIVATE")
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
	postgresUrl             string
	postgresDbName          string
	postgresUser            string
	postgresPassword        string
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
func (env Env) PostgresUrl() string {
	return env.postgresUrl
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

func (devEnv DevEnv) PostgresDbName() string {
	return devEnv.postgresDbName
}
func (devEnv DevEnv) PostgresUser() string {
	return devEnv.postgresUser
}
func (devEnv DevEnv) PostgresPassword() string {
	return devEnv.postgresPassword
}

/**
* Loads from env by key or panics
 */
func loadFromEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("env: %v needs to be defined", key)
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
func (stagingEnv StagingEnv) PostgresDbName() string {
	return stagingEnv.postgresDbName
}
func (stagingEnv StagingEnv) PostgresUser() string {
	return stagingEnv.postgresUser
}
func (stagingEnv StagingEnv) PostgresPassword() string {
	return stagingEnv.postgresPassword
}

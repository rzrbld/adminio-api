package config

import (
	"os"
	strconv "strconv"
)

var (
	Server         = getEnv("MINIO_HOST_PORT", "localhost:9000")
	Maccess        = getEnv("MINIO_ACCESS", "test")
	Msecret        = getEnv("MINIO_SECRET", "testtest123")
	Region         = getEnv("MINIO_REGION", "us-east-1")
	Ssl, _         = strconv.ParseBool(getEnv("MINIO_SSL", "false"))
	ServerHostPort = getEnv("ADMINIO_HOST_PORT", "localhost:8080")
	AdminioCORS    = getEnv("ADMINIO_CORS_DOMAIN", "*")
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	ScHashKey  = getEnv("ADMINIO_COOKIE_HASH_KEY", "NRUeuq6AdskNPa7ewZuxG9TrDZC4xFat")
	ScBlockKey = getEnv("ADMINIO_COOKIE_BLOCK_KEY", "bnfYuphzxPhJMR823YNezH83fuHuddFC")
	// ---------------
	ScCookieName      = getEnv("ADMINIO_COOKIE_NAME", "adminiosessionid")
	OauthEnable, _    = strconv.ParseBool(getEnv("ADMINIO_OAUTH_ENABLE", "false"))
	AuditLogEnable, _ = strconv.ParseBool(getEnv("ADMINIO_AUDIT_LOG_ENABLE", "false"))
	MetricsEnable, _  = strconv.ParseBool(getEnv("ADMINIO_METRICS_ENABLE", "false"))
	OauthProvider     = getEnv("ADMINIO_OAUTH_PROVIDER", "github")
	OauthClientId     = getEnv("ADMINIO_OAUTH_CLIENT_ID", "my-github-oauth-app-client-id")
	OauthClientSecret = getEnv("ADMINIO_OAUTH_CLIENT_SECRET", "my-github-oauth-app-secret")
	OauthCallback     = getEnv("ADMINIO_OAUTH_CALLBACK", "http://"+ServerHostPort+"/auth/callback")
	OauthCustomDomain = getEnv("ADMINIO_OAUTH_CUSTOM_DOMAIN", "")
)

func getEnv(key, fallback string) string {
	value, exist := os.LookupEnv(key)

	if !exist {
		return fallback
	}

	return value
}

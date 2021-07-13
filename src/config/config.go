package config

import (
	"os"
	strconv "strconv"

	log "github.com/sirupsen/logrus"
)

var (
	Server      = getEnv("MINIO_HOST_PORT", "localhost:9000")
	Maccess     = getEnv("MINIO_ACCESS", "test")
	Msecret     = getEnv("MINIO_SECRET", "testtest123")
	Region      = getEnv("MINIO_REGION", "us-east-1")
	SvcLogLevel = getEnv("ADMINIO_LOG_LEVEL", "INFO")
	// Enable object locking by default
	DefaultObjectLocking, _ = strconv.ParseBool(getEnv("MINIO_DEFAULT_LOCK_OBLECT_ENABLE", "false"))
	Ssl, _                  = strconv.ParseBool(getEnv("MINIO_SSL", "false"))
	ServerHostPort          = getEnv("ADMINIO_HOST_PORT", "localhost:8080")
	AdminioCORS             = getEnv("ADMINIO_CORS_DOMAIN", "*")
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	ScHashKey  = getEnv("ADMINIO_COOKIE_HASH_KEY", "NRUeuq6AdskNPa7ewZuxG9TrDZC4xFat")
	ScBlockKey = getEnv("ADMINIO_COOKIE_BLOCK_KEY", "bnfYuphzxPhJMR823YNezH83fuHuddFC")
	// ---------------
	ScCookieName      = getEnv("ADMINIO_COOKIE_NAME", "adminiosessionid")
	OauthEnable, _    = strconv.ParseBool(getEnv("ADMINIO_OAUTH_ENABLE", "false"))
	AuditLogEnable, _ = strconv.ParseBool(getEnv("ADMINIO_AUDIT_LOG_ENABLE", "false"))
	MetricsEnable, _  = strconv.ParseBool(getEnv("ADMINIO_METRICS_ENABLE", "false"))
	ProbesEnable, _   = strconv.ParseBool(getEnv("ADMINIO_PROBES_ENABLE", "false"))
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

func SetLogLevel() {
	log.Infoln("Set log level to: ", SvcLogLevel)
	selectedLogLevel := log.InfoLevel
	switch loglvl := SvcLogLevel; loglvl {
	case "TRACE":
		selectedLogLevel = log.TraceLevel
	case "DEBUG":
		selectedLogLevel = log.DebugLevel
	case "INFO":
		selectedLogLevel = log.InfoLevel
	case "WARN":
		selectedLogLevel = log.WarnLevel
	case "ERROR":
		selectedLogLevel = log.ErrorLevel
	case "FATAL":
		selectedLogLevel = log.FatalLevel
	case "PANIC":
		selectedLogLevel = log.PanicLevel
	default:
		log.Errorln("Unknown log level:", SvcLogLevel, ". Fallback to INFO.", "Possible values: \n TRACE \n DEBUG \n INFO \n WARN \n ERROR \n FATAL \n PANIC \n")
	}

	log.SetLevel(selectedLogLevel)
}

package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/G0tem/go-servise-auth/internal"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel                           int      `default:"4" envconfig:"LOG_LEVEL"`
	HttpPort                           uint16   `default:"8002" envconfig:"HTTP_PORT"`
	UserServiceBaseUrl                 string   `binding:"required" envconfig:"USER_SERVICE_BASE_URL"`
	CorsOrigins                        []string `binding:"required" envconfig:"CORS_ORIGINS"`
	SecretKey                          string   `binding:"required" envconfig:"SECRET_KEY"`
	PublicEmailConfirmationUrl         string   `binding:"required" envconfig:"PUBLIC_EMAIL_CONFIRMATION_URL"`
	PublicPasswordResetConfirmationUrl string   `binding:"required" envconfig:"PUBLIC_PASSWORD_RESET_CONFIRMATION_URL"`
	PublicUrl                          string   `binding:"required" envconfig:"PUBLIC_URL"`
	PublicErrorUrl                     string   `binding:"required" envconfig:"PUBLIC_ERROR_URL"`

	PostgresHost            string        `binding:"required" envconfig:"POSTGRES_HOST"`
	PostgresPort            string        `binding:"required" envconfig:"POSTGRES_PORT"`
	PostgresDb              string        `binding:"required" envconfig:"POSTGRES_DB"`
	PostgresUser            string        `binding:"required" envconfig:"POSTGRES_USER"`
	PostgresPassword        string        `binding:"required" envconfig:"POSTGRES_PASSWORD"`
	PostgresMaxIdleConns    int           `default:"10" envconfig:"POSTGRES_MAX_IDLE_CONNS"`
	PostgresMaxOpenConns    int           `default:"100" envconfig:"POSTGRES_MAX_OPEN_CONNS"`
	PostgresConnMaxLifetime time.Duration `default:"1h" envconfig:"POSTGRES_CONN_MAX_LIFETIME"`

	JwtValidationUrl string `binding:"required" envconfig:"JWT_VALIDATION_URL"`

	RMQConnUrl                string `binding:"required" envconfig:"RMQ_CONN_URL"`
	RMQMailExchange           string `binding:"required" envconfig:"RMQ_MAIL_EXCHANGE"`
	RMQMailExchangeAutocreate bool   `binding:"required" envconfig:"RMQ_MAIL_EXCHANGE_AUTOCREATE_ENABLED"`

	RedisAddr string `binding:"required" envconfig:"REDIS_ADDR"`
	RedisDB   int    `binding:"required" envconfig:"REDIS_DB"`

	CdnPublicUrl        string `binding:"required" envconfig:"CDN_PUBLIC_URL"`
	S3AvatarsBucketName string `binding:"required" envconfig:"S3_AVATARS_BUCKET_NAME"`
	S3CoversBucketName  string `binding:"required" envconfig:"S3_COVERS_BUCKET_NAME"`
	S3Region            string `binding:"required" envconfig:"S3_REGION"`
	S3Endpoint          string `binding:"required" envconfig:"S3_ENDPOINT"`
	S3AccessKey         string `binding:"required" envconfig:"S3_ACCESS_KEY"`
	S3SecretAccessKey   string `binding:"required" envconfig:"S3_SECRET_ACCESS_KEY"`

	MaxFileUploadSizeInBytes int `default:"10485760" envconfig:"MAX_FILE_UPLOAD_SIZE"`
}

func getenvDef(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Debug().Msg("use process environment variables (i.e. not from .env)")
	}

	logLevel, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))

	return Config{
		LogLevel:                           logLevel,
		HttpPort:                           internal.ParseUint16(os.Getenv("HTTP_PORT"), 8002),
		UserServiceBaseUrl:                 getenvDef("USER_SERVICE_BASE_URL", "http://31.131.255.218:8080/api"),
		CorsOrigins:                        strings.Split(os.Getenv("CORS_ORIGINS"), ","),
		SecretKey:                          os.Getenv("SECRET_KEY"),
		PublicEmailConfirmationUrl:         os.Getenv("PUBLIC_EMAIL_CONFIRMATION_URL"),
		PublicPasswordResetConfirmationUrl: os.Getenv("PUBLIC_PASSWORD_RESET_CONFIRMATION_URL"),
		PublicUrl:                          os.Getenv("PUBLIC_URL"),
		PublicErrorUrl:                     os.Getenv("PUBLIC_ERROR_URL"),

		PostgresHost:            os.Getenv("POSTGRES_HOST"),
		PostgresPort:            os.Getenv("POSTGRES_PORT"),
		PostgresDb:              os.Getenv("POSTGRES_DB"),
		PostgresUser:            os.Getenv("POSTGRES_USER"),
		PostgresPassword:        os.Getenv("POSTGRES_PASSWORD"),
		PostgresMaxIdleConns:    internal.ParseInt(os.Getenv("POSTGRES_MAX_IDLE_CONNS"), 10),
		PostgresMaxOpenConns:    internal.ParseInt(os.Getenv("POSTGRES_MAX_OPEN_CONNS"), 100),
		PostgresConnMaxLifetime: internal.ParseDuration(os.Getenv("POSTGRES_CONN_MAX_LIFETIME"), 1*time.Hour),

		JwtValidationUrl: os.Getenv("JWT_VALIDATION_URL"),

		RMQConnUrl:                os.Getenv("RMQ_CONN_URL"),
		RMQMailExchange:           os.Getenv("RMQ_MAIL_EXCHANGE"),
		RMQMailExchangeAutocreate: internal.ParseBool(os.Getenv("RMQ_MAIL_EXCHANGE_AUTOCREATE_ENABLED")),

		RedisAddr: os.Getenv("REDIS_ADDR"),
		RedisDB:   internal.ParseInt(os.Getenv("REDIS_DB"), 0),

		CdnPublicUrl:        os.Getenv("CDN_PUBLIC_URL"),
		S3AvatarsBucketName: os.Getenv("S3_AVATARS_BUCKET_NAME"),
		S3CoversBucketName:  os.Getenv("S3_COVERS_BUCKET_NAME"),
		S3Region:            os.Getenv("S3_REGION"),
		S3Endpoint:          os.Getenv("S3_ENDPOINT"),
		S3AccessKey:         os.Getenv("S3_ACCESS_KEY"),
		S3SecretAccessKey:   os.Getenv("S3_SECRET_ACCESS_KEY"),

		MaxFileUploadSizeInBytes: internal.ParseInt(os.Getenv("MAX_FILE_UPLOAD_SIZE"), 10485760),
	}
}

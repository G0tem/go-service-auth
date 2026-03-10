package config

import (
	"os"
	"strconv"
	"time"

	"github.com/G0tem/go-service-auth/internal"
)

type Config struct {
	LogLevel int `default:"4" envconfig:"LOG_LEVEL"`

	HttpPort uint16 `default:"8002" envconfig:"HTTP_PORT"`
	GrpcPort uint16 `default:"50051" envconfig:"GRPC_PORT"`

	SecretKey string `binding:"required" envconfig:"SECRET_KEY"`

	PostgresHost            string        `binding:"required" envconfig:"POSTGRES_HOST"`
	PostgresPort            string        `binding:"required" envconfig:"POSTGRES_PORT"`
	PostgresDb              string        `binding:"required" envconfig:"POSTGRES_DB"`
	PostgresUser            string        `binding:"required" envconfig:"POSTGRES_USER"`
	PostgresPassword        string        `binding:"required" envconfig:"POSTGRES_PASSWORD"`
	PostgresMaxIdleConns    int           `default:"10" envconfig:"POSTGRES_MAX_IDLE_CONNS"`
	PostgresMaxOpenConns    int           `default:"100" envconfig:"POSTGRES_MAX_OPEN_CONNS"`
	PostgresConnMaxLifetime time.Duration `default:"1h" envconfig:"POSTGRES_CONN_MAX_LIFETIME"`

	RMQConnUrl                string `binding:"required" envconfig:"RMQ_CONN_URL"`
	RMQMailExchange           string `binding:"required" envconfig:"RMQ_MAIL_EXCHANGE"`
	RMQMailExchangeAutocreate bool   `binding:"required" envconfig:"RMQ_MAIL_EXCHANGE_AUTOCREATE_ENABLED"`

	RedisAddr string `binding:"required" envconfig:"REDIS_ADDR"`
	RedisDB   int    `binding:"required" envconfig:"REDIS_DB"`

	S3AvatarsBucketName      string `binding:"required" envconfig:"S3_AVATARS_BUCKET_NAME"`
	S3CoversBucketName       string `binding:"required" envconfig:"S3_COVERS_BUCKET_NAME"`
	S3Region                 string `binding:"required" envconfig:"S3_REGION"`
	S3Endpoint               string `binding:"required" envconfig:"S3_ENDPOINT"`
	S3AccessKey              string `binding:"required" envconfig:"S3_ACCESS_KEY"`
	S3SecretAccessKey        string `binding:"required" envconfig:"S3_SECRET_ACCESS_KEY"`
	MaxFileUploadSizeInBytes int    `default:"10485760" envconfig:"MAX_FILE_UPLOAD_SIZE"`
}

func LoadConfig() Config {
	logLevel, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))

	return Config{
		LogLevel: logLevel,

		HttpPort: internal.ParseUint16(os.Getenv("HTTP_PORT"), 8002),
		GrpcPort: internal.ParseUint16(os.Getenv("GRPC_PORT"), 50051),

		SecretKey: os.Getenv("SECRET_KEY"),

		PostgresHost:            os.Getenv("POSTGRES_HOST"),
		PostgresPort:            os.Getenv("POSTGRES_PORT"),
		PostgresDb:              os.Getenv("POSTGRES_DB"),
		PostgresUser:            os.Getenv("POSTGRES_USER"),
		PostgresPassword:        os.Getenv("POSTGRES_PASSWORD"),
		PostgresMaxIdleConns:    internal.ParseInt(os.Getenv("POSTGRES_MAX_IDLE_CONNS"), 10),
		PostgresMaxOpenConns:    internal.ParseInt(os.Getenv("POSTGRES_MAX_OPEN_CONNS"), 100),
		PostgresConnMaxLifetime: internal.ParseDuration(os.Getenv("POSTGRES_CONN_MAX_LIFETIME"), 1*time.Hour),

		RMQConnUrl:                os.Getenv("RMQ_CONN_URL"),
		RMQMailExchange:           os.Getenv("RMQ_MAIL_EXCHANGE"),
		RMQMailExchangeAutocreate: internal.ParseBool(os.Getenv("RMQ_MAIL_EXCHANGE_AUTOCREATE_ENABLED")),

		RedisAddr: os.Getenv("REDIS_ADDR"),
		RedisDB:   internal.ParseInt(os.Getenv("REDIS_DB"), 0),

		S3AvatarsBucketName:      os.Getenv("S3_AVATARS_BUCKET_NAME"),
		S3CoversBucketName:       os.Getenv("S3_COVERS_BUCKET_NAME"),
		S3Region:                 os.Getenv("S3_REGION"),
		S3Endpoint:               os.Getenv("S3_ENDPOINT"),
		S3AccessKey:              os.Getenv("S3_ACCESS_KEY"),
		S3SecretAccessKey:        os.Getenv("S3_SECRET_ACCESS_KEY"),
		MaxFileUploadSizeInBytes: internal.ParseInt(os.Getenv("MAX_FILE_UPLOAD_SIZE"), 10485760),
	}
}

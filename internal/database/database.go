package database

import (
	"fmt"
	"strings"

	"github.com/G0tem/go-service-auth/internal"
	"github.com/G0tem/go-service-auth/internal/config"
	"github.com/G0tem/go-service-auth/internal/model"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Connect function
func Connect(cfg config.Config) (*gorm.DB, error) {
	postgresPort := cfg.PostgresPort
	// because our config function returns a string, we are parsing our      str to int here
	port := internal.ParseInt(postgresPort, 5432)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDb,
		port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		Logger: &internal.GormZeroLogAdapter{
			Level: cfg.LogLevel,
		},
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",                                // убираем префикс
			SingularTable: false,                             // use plural table name, table for `User` would be `users`
			NoLowerCase:   false,                             // skip the snake_casing of names
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
	})

	if err != nil {
		log.Error().Msgf("failed to connect to database. %v\n", err)
		return nil, err
	}

	log.Info().Msg("Connected")
	log.Info().Msg("running migrations")
	err = db.AutoMigrate(
		&model.User{},
		&model.UserRole{},
		&model.UserPermission{},
		&model.UserRolePermission{},
	)
	if err != nil {
		log.Error().Msgf("failed run auto-migrations. %v\n", err)
		return nil, err
	}

	// Apply connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Msgf("failed to create database connection pool. %v\n", err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.PostgresMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.PostgresMaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.PostgresConnMaxLifetime)

	return db, nil
}

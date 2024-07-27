package config

import "e-resep-be/internal/helper"

type (
	Configuration struct {
		Server   *Server
		Database *Database
	}

	Server struct {
		AppPort int
		AppEnv  string
		AppName string
		AppID   string
	}

	Database struct {
		Port     int
		Host     string
		Username string
		Password string
		Name     string
		SslMode  string
	}
)

func loadConfiguration() *Configuration {
	return &Configuration{
		Server: &Server{
			AppPort: helper.GetEnvInt("APP_PORT"),
			AppEnv:  helper.GetEnvString("APP_ENV"),
			AppName: helper.GetEnvString("APP_NAME"),
			AppID:   helper.GetEnvString("APP_ID"),
		},
		Database: &Database{
			Port:     helper.GetEnvInt("DB_PORT"),
			Host:     helper.GetEnvString("DB_HOST"),
			Username: helper.GetEnvString("DB_USERNAME"),
			Password: helper.GetEnvString("DB_PASSWORD"),
			Name:     helper.GetEnvString("DB_NAME"),
			SslMode:  helper.GetEnvString("DB_SSL_MODE"),
		},
	}
}

func NewConfig() *Configuration {
	return loadConfiguration()
}

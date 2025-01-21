package config

import "os"

type Environment struct {
	DATABASE_HOST     string
	DATABASE_USER     string
	DATABASE_PASSWORD string
	DATABASE_NAME     string
	DATABASE_PORT     string
	DATABASE_SCHEMA   string
	JWT_SECRET_KEY    string
	PORT              string
}

func EnvironmentConfig() Environment {
	return Environment{
		DATABASE_HOST:     os.Getenv("DATABASE_HOST"),
		DATABASE_USER:     os.Getenv("DATABASE_USER"),
		DATABASE_PASSWORD: os.Getenv("DATABASE_PASSWORD"),
		DATABASE_NAME:     os.Getenv("DATABASE_NAME"),
		DATABASE_PORT:     os.Getenv("DATABASE_PORT"),
		DATABASE_SCHEMA:   os.Getenv("DATABASE_SCHEMA"),
		JWT_SECRET_KEY:    os.Getenv("JWT_SECRET_KEY"),
		PORT:              os.Getenv("PORT"),
	}
}

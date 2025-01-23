package config

import "os"

type Environment struct {
	DATABASE_HOST         string
	DATABASE_USER         string
	DATABASE_PASSWORD     string
	DATABASE_NAME         string
	DATABASE_PORT         string
	DATABASE_SCHEMA       string
	JWT_SECRET            string
	PORT                  string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_S3_BUCKET         string
	AWS_S3_REGION         string
}

func EnvironmentConfig() Environment {
	return Environment{
		DATABASE_HOST:         os.Getenv("DATABASE_HOST"),
		DATABASE_USER:         os.Getenv("DATABASE_USER"),
		DATABASE_PASSWORD:     os.Getenv("DATABASE_PASSWORD"),
		DATABASE_NAME:         os.Getenv("DATABASE_NAME"),
		DATABASE_PORT:         os.Getenv("DATABASE_PORT"),
		DATABASE_SCHEMA:       os.Getenv("DATABASE_SCHEMA"),
		JWT_SECRET:            os.Getenv("JWT_SECRET"),
		PORT:                  os.Getenv("PORT"),
		AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWS_S3_BUCKET:         os.Getenv("AWS_S3_BUCKET"),
		AWS_S3_REGION:         os.Getenv("AWS_S3_REGION"),
	}
}

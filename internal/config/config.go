package config

import "gemsvietnambe/internal/model/database"

type ApiConfig struct {
	Env     string
	Port    string
	Version string
	DB      *database.Queries
}

const (
	DEV_ENV   = "DEV"
	STAGE_ENV = "STAGE"
	PROD_ENV  = "PROD"
)

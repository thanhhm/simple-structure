package main

import (
	"log"
	"os"

	"github.com/subosito/gotenv"
)

type config struct {
	ginPort             string
	mysqlDataSourceName string
}

func newConfig() *config {
	if err := gotenv.Load(); err != nil {
		log.Fatal("Missing environment file")
	}

	return &config{
		ginPort:             os.Getenv("GIN_PORT"),
		mysqlDataSourceName: os.Getenv("MYSQL_DATA_SOURCE_NAME"),
	}
}

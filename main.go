package main

import (
	"log"
	"os"

	"github.com/stinkyfingers/stageplotapi/server"
	"github.com/stinkyfingers/stageplotapi/storage"
)

func main() {
	cfg, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	s := server.NewServer(cfg.Port, cfg.Storage)
	log.Fatal(s.Run())
}

// config
type Config struct {
	Storage storage.Storage
	Port    string
}

func GetConfig() (*Config, error) {
	store, err := storage.NewMongo()
	if err != nil {
		return nil, err
	}
	port := "8080"
	if val := os.Getenv("PORT"); val != "" {
		port = val
	}

	return &Config{
		Storage: store,
		Port:    port,
	}, nil
}

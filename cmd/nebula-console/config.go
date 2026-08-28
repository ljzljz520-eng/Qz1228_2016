package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabasePath    string
	Address         string
	ShutdownTimeout time.Duration
	Seed            bool
}

func LoadConfig() Config {
	cfg := Config{DatabasePath: "nebula.db", Address: ":8080", ShutdownTimeout: 5 * time.Second, Seed: true}
	if value := os.Getenv("NEBULA_DB"); value != "" {
		cfg.DatabasePath = value
	}
	if value := os.Getenv("NEBULA_ADDR"); value != "" {
		cfg.Address = value
	}
	if value := os.Getenv("NEBULA_SEED"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			cfg.Seed = parsed
		}
	}
	return cfg
}

func (c Config) Validate() error {
	if c.DatabasePath == "" || c.Address == "" || c.ShutdownTimeout <= 0 {
		return fmt.Errorf("invalid configuration")
	}
	return nil
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

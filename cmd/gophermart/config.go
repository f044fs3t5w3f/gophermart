package main

import (
	"flag"
	"os"
)

type Config struct {
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
}

func getConfig() *Config {
	runAddressFlag := flag.String("a", "localhost:8080", "server address")
	databaseURIFlag := flag.String("d", "", "database uri")
	accrualSystemAddressFlag := flag.String("r", "", "accrual System Address")

	flag.Parse()
	config := &Config{}
	config.runAddress = envOrString("RUN_ADDRESS", *runAddressFlag)
	config.databaseURI = envOrString("DATABASE_URI", *databaseURIFlag)
	config.accrualSystemAddress = envOrString("ACCRUAL_SYSTEM_ADDRESS", *accrualSystemAddressFlag)
	return config
}

func envOrString(env string, fallback string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return fallback
}

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Lyckabc/ipTIME-API/cmd/structs"
)

func RouterInfo() *structs.Router {
	loadEnv()

	port := 80
	if p := os.Getenv("IPTIME_PORT"); p != "" {
		parsed, err := strconv.Atoi(p)
		if err != nil {
			log.Fatalf("invalid IPTIME_PORT: %s", p)
		}
		port = parsed
	}

	host := os.Getenv("IPTIME_HOST")
	if host == "" {
		log.Fatal("IPTIME_HOST is required")
	}

	return &structs.Router{
		Host:     host,
		Port:     port,
		Username: os.Getenv("IPTIME_USERNAME"),
		Password: os.Getenv("IPTIME_PASSWORD"),
	}
}

// loadEnv reads a .env file, searching the working directory and up to two parent directories.
func loadEnv() {
	candidates := []string{".env", "../.env", "../../.env"}
	var data []byte
	var err error
	for _, path := range candidates {
		data, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if os.Getenv(key) == "" {
			os.Setenv(key, fmt.Sprintf("%s", val))
		}
	}
}

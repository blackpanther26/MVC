package config

import (
	"log"
	"os"
	"path/filepath"
	"fmt"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	getEnvFilePath := func(envFileName string) (string, error) {
		currentDir, err := os.Getwd()
		if err!= nil {
			return "", fmt.Errorf("error getting current directory: %w", err)
		}

		for {
			goModPath := filepath.Join(currentDir, "go.mod")
			if _, err := os.Stat(goModPath); err == nil {
				envFilePath := filepath.Join(currentDir, envFileName)
				return envFilePath, nil
			}

			parent := filepath.Dir(currentDir)
			if parent == currentDir {
				return "", fmt.Errorf("go.mod not found")
			}

			currentDir = parent
		}
	}

	envFilePath, err := getEnvFilePath(".env")
	if err!= nil {
		log.Fatalf("Error locating.env file: %v", err)
	}

	err = godotenv.Load(envFilePath)
	if err!= nil {
		log.Fatalf("Error loading.env file: %v", err)
	}
}

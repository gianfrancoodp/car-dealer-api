package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// This function checks if the environment variable is correctly loaded, and if it exist, returns the variable.
func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!!")
	}

	return os.Getenv("MONGOURI")
}

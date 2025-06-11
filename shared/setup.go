package shared

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

const APP_ENV = "APP_ENV"

func setupEnv() string {
	env := os.Getenv(APP_ENV)

	if env == "development" {
		log.SetLevel(log.DebugLevel)
	} else if env == "production" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.Info("Unknown log level set, setting development")
		env = "development"
		log.SetLevel(log.DebugLevel)
	}

	log.Info(fmt.Sprintf("Running in %s", env))

	return env
}

func initializeEnvFile() {
	err := godotenv.Load("../.env")
	// Load environment variables from .env file
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
}

func initializeSupabase() *supabase.Client {
	supabaseClient, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_ANON_KEY"), &supabase.ClientOptions{})

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return supabaseClient
}

func Setup() (*supabase.Client, string) {
	env := setupEnv()
	initializeEnvFile()
	supabaseClient := initializeSupabase()

	return supabaseClient, env
}

package shared

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

const APP_ENV = "APP_ENV"

type Env string

const (
	PROD  Env = "prod"
	STAGE Env = "stage"
	DEV   Env = "dev"
)

func SetupEnv() Env {
	env := Env(os.Getenv(APP_ENV))

	if env == DEV || env == STAGE {
		log.SetLevel(log.DebugLevel)
	} else if env == PROD {
		log.SetLevel(log.InfoLevel)
	} else {
		log.Info("Unknown log level set, setting development")
		env = DEV
		log.SetLevel(log.DebugLevel)
	}

	log.Info(fmt.Sprintf("We in %s baby", env))

	return env
}

func initializeEnvFile(filePath string) {
	env := SetupEnv()
	if env == "prod" || env == "stage" {
		log.Debug("Loading variables via docker container")
		return
	}

	var err error
	if filePath == "" {
		err = godotenv.Load(".env")
	} else {
		err = godotenv.Load(filePath)
	}

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

func Setup(filePath string) (*supabase.Client, Env) {
	env := SetupEnv()
	initializeEnvFile(filePath)
	supabaseClient := initializeSupabase()

	return supabaseClient, env
}

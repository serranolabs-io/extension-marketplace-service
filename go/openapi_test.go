package openapi

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func TestFile(t *testing.T) {

	directory := getDirectory("new-extension", "FfjugE")
	packgeJson := readEndpoint(path.Join(directory, PACKAGE_JSON), false)
	configJson := readEndpoint(path.Join(directory, CONFIG_JSON), false)

	fmt.Println(packgeJson)
	fmt.Println(configJson)
}

func TestAll(t *testing.T) {
	directory := getFullDirectory()
	packgeJson := readEndpoint(directory, true)
	configJson := readEndpoint(directory, true)

	fmt.Println(packgeJson)
	fmt.Println(configJson)
}

func initializeSupabase() *supabase.Client {
	supabaseClient, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_ANON_KEY"), &supabase.ClientOptions{})

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return supabaseClient
}
func initializeEnvVariables() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
}

func TestGetAllExtensions(t *testing.T) {
	initializeEnvVariables()
	initializeSupabase()
	supabaseClient := initializeSupabase()

	getAllExtensions(supabaseClient)
}

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/api/idtoken"
)

var (
	db *sql.DB
)

func main() {
	var err error

	databaseURL := os.Getenv("DATABASE_URL")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")

	
	fmt.Printf("Connecting to database: %s\n", databaseURL)
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	
	defer db.Close()
	
	createTable()

	http.HandleFunc("/auth/google", googleAuthHandler)
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		google_id VARCHAR(255) UNIQUE NOT NULL
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Error creating table: %q", err)
	}
}

func googleAuthHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDToken string `json:"idToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tokenInfo, err := validateGoogleIDToken(req.IDToken)
	if err != nil {
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}

	googleID := tokenInfo.Subject


	if err := upsertUser(googleID); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User authenticated")
}

func validateGoogleIDToken(idToken string) (*idtoken.Payload, error) {
	ctx := context.Background()
	validator, err := idtoken.NewValidator(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create token validator: %w", err)
	}

	payload, err := validator.Validate(ctx, idToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, fmt.Errorf("failed to validate ID token: %w", err)
	}

	return payload, nil
}

func upsertUser(googleID string) error {
	query := `
	INSERT INTO users (google_id)
	VALUES ($1)
	ON CONFLICT (google_id)
	DO NOTHING;
	`
	_, err := db.Exec(query, googleID)
	return err
}

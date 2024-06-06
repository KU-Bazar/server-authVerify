// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"
)


func main() {

	http.HandleFunc("/auth/google", googleAuthHandler)
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
	fmt.Println(googleID)

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


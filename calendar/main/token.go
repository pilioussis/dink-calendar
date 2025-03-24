package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

const DEAN_TOKEN = "dean_token.json"
const STRUGS_TOKEN = "strugs_token.json"

const CREDENTIALS_FOLDER = "credentials"

const CREDENTIALS_FILE = "credentials.json"

func getClient(config *oauth2.Config, token string) *http.Client {
	tokenFile := fmt.Sprintf("%s/%s", CREDENTIALS_FOLDER, token)
	tok, err := tokenFromFile(tokenFile)
	if err != nil || !tok.Valid() {
		if err != nil {
			fmt.Println("Error reading token file:", err)
		} else {
			fmt.Println("Token is invalid, getting new token")
		}
		tok = getTokenFromWeb(config)
		saveTokenToFile(tokenFile, tok)
	} else if tok.Expiry.After(time.Now().Add(-12 * time.Hour)) {
		tok = refreshToken(config, tok, tokenFile)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser: \n%v\n\n", authURL)
	fmt.Println("Copy the `code` query param of the URL you are redirected to and enter it here")
	fmt.Println("\nENTER CODE HERE:")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read auth code: %v", err)
	}
	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func getToken(token string) *http.Client {
	credentials, err := os.ReadFile(fmt.Sprintf("%s/%s", CREDENTIALS_FOLDER, CREDENTIALS_FILE))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(credentials, calendar.CalendarReadonlyScope)

	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, token)

	return client
}

func refreshToken(config *oauth2.Config, tok *oauth2.Token, tokenFile string) *oauth2.Token {
	fmt.Println("Refreshing", tokenFile)
	ctx := context.Background()

	// Create a token object with the refresh token
	token := &oauth2.Token{
		RefreshToken: tok.RefreshToken,
	}

	// Use the token source to refresh the token
	tokenSource := config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Panicln("Error refreshing token:", err)
	}

	// Save the new token to a file
	err = saveTokenToFile(tokenFile, newToken)
	if err != nil {
		log.Panicln("Error saving token to file:", err)
	}

	fmt.Println("Refreshed")
	return newToken
}

// saveTokenToFile saves the token to a file in JSON format
func saveTokenToFile(filename string, token *oauth2.Token) error {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling token: %w", err)
	}

	err = os.WriteFile(filename, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing token to file: %w", err)
	}

	return nil
}

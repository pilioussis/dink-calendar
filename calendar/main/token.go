package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

const OAUTH_TOKEN_FILE = "token.json"
const CREDENTIALS_FOLDER = "credentials"
const CREDENTIALS_FILE = "credentials.json"

func getClient(token string) *http.Client {
	config := getConfig(token)
	tokenFile := fmt.Sprintf("%s/%s", CREDENTIALS_FOLDER, token)
	tok, err := tokenFromFile(tokenFile)

	if err == nil {
		tok, err = refreshToken(config, tok, tokenFile)
		if err == nil {
			return config.Client(context.Background(), tok)
		}
		slog.Warn("Token refresh unavailable", "error", err)
	} else {
		slog.Warn("No token file", "error", err)
	}

	tok = getTokenFromWeb(config)
	saveTokenToFile(tokenFile, tok)

	return config.Client(context.Background(), tok)
}

func getConfig(token string) *oauth2.Config {
	credentials, err := os.ReadFile(fmt.Sprintf("%s/%s", CREDENTIALS_FOLDER, CREDENTIALS_FILE))
	if err != nil {
		slog.Error("Unable to read client secret file", "error", err)
		panic(err)
	}

	config, err := google.ConfigFromJSON(credentials, calendar.CalendarReadonlyScope)
	if err != nil {
		slog.Error("Unable to parse oauth config file to config", "error", err)
		panic(err)
	}

	return config
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
	slog.Info("Requires manual user auth via web")
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("\nGo to the following link in your browser: \n%s\n\n", authURL)
	fmt.Println("Copy the `code` query param of the URL you are redirected to and enter it here")
	fmt.Println("\nENTER CODE HERE:")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		slog.Error("Unable to read auth code", "error", err)
		panic(err)
	}
	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		slog.Error("Unable to retrieve token from web", "error", err)
		panic(err)
	}
	return tok
}

func refreshToken(config *oauth2.Config, tok *oauth2.Token, tokenFile string) (*oauth2.Token, error) {
	slog.Info("Refreshing", "tokenFile", tokenFile)
	ctx := context.Background()

	tokenSource := config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: tok.RefreshToken,
	})
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	err = saveTokenToFile(tokenFile, newToken)
	if err != nil {
		slog.Error("Error saving token to file", "error", err)
		panic(err)
	}

	slog.Info("Refreshed token", "expiry", newToken.Expiry)
	return newToken, nil
}

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

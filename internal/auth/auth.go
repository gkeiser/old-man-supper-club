package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Config holds the OAuth2 configuration and the whitelist.
type Config struct {
	OAuth         *oauth2.Config
	AllowedEmails map[string]bool
}

// NewConfig initializes the OAuth2 settings and the email whitelist.
func NewConfig() *Config {
	oauth := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	allowed := make(map[string]bool)
	emails := strings.Split(os.Getenv("ALLOWED_EMAILS"), ",")
	for _, e := range emails {
		allowed[strings.TrimSpace(e)] = true
	}

	return &Config{
		OAuth:         oauth,
		AllowedEmails: allowed,
	}
}

// GoogleUser represents the data returned by Google's userinfo API.
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GetLoginURL generates the Google login link.
func (c *Config) GetLoginURL(state string) string {
	return c.OAuth.AuthCodeURL(state)
}

// VerifyUser exchanges the code for user info and checks the whitelist.
func (c *Config) VerifyUser(ctx context.Context, code string) (*GoogleUser, error) {
	token, err := c.OAuth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	client := c.OAuth.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// CHECK WHITELIST
	if !c.AllowedEmails[user.Email] {
		return nil, fmt.Errorf("unauthorized: email %s is not in the whitelist", user.Email)
	}

	return &user, nil
}

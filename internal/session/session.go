package session

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var (
	// store uses a secret key to sign the cookie.
	store *sessions.CookieStore
	name  = "supper-club-session"
)

// Init initializes the session store.
func Init() {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "super-secret-fallback-key"
	}
	store = sessions.NewCookieStore([]byte(secret))

	// Security settings
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production", // Only send over HTTPS in prod
		SameSite: http.SameSiteLaxMode,
	}
}

// SetUser saves the user ID and name in the session.
func SetUser(w http.ResponseWriter, r *http.Request, userID, userName string) error {
	session, _ := store.Get(r, name)
	session.Values["user_id"] = userID
	session.Values["user_name"] = userName
	return session.Save(r, w)
}

// GetUser retrieves the user ID and name from the session.
func GetUser(r *http.Request) (string, string) {
	session, _ := store.Get(r, name)
	userID, _ := session.Values["user_id"].(string)
	userName, _ := session.Values["user_name"].(string)
	return userID, userName
}

// Logout clears the session.
func Logout(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, name)
	session.Values["user_id"] = ""
	session.Values["user_name"] = ""
	session.Options.MaxAge = -1 // Expire the cookie immediately
	return session.Save(r, w)
}

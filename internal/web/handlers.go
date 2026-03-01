package web

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/grant/old-man-supper-club/internal/auth"
	"github.com/grant/old-man-supper-club/internal/db"
	"github.com/grant/old-man-supper-club/internal/models"
	"github.com/grant/old-man-supper-club/internal/scoring"
	"github.com/grant/old-man-supper-club/internal/session"
)

type Server struct {
	AuthConfig *auth.Config
	Repo       *db.Repository
	Templates  map[string]*template.Template
}

type TemplateData struct {
	IsAuthenticated bool
	UserName        string
	Restaurants     []models.Restaurant
	Restaurant      *models.Restaurant
	Reviews         []models.Review
	OverallScore    float64
	Weights         map[string]float64
}

func NewServer(authCfg *auth.Config, repo *db.Repository, tmplFS embed.FS) *Server {
	templates := make(map[string]*template.Template)

	funcMap := template.FuncMap{
		"multiply": func(a, b float64) float64 {
			return a * b
		},
	}

	pages := []string{"home.html", "leaderboard.html", "restaurant_details.html", "add_restaurant.html"}
	for _, page := range pages {
		tmpl := template.Must(template.New(page).Funcs(funcMap).ParseFS(tmplFS, "templates/base.html", "templates/"+page))
		templates[page] = tmpl
	}

	return &Server{
		AuthConfig: authCfg,
		Repo:       repo,
		Templates:  templates,
	}
}

func (s *Server) render(w http.ResponseWriter, name string, data TemplateData) {
	tmpl, ok := s.Templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleHome shows the landing page.
func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	userID, userName := session.GetUser(r)
	s.render(w, "home.html", TemplateData{IsAuthenticated: userID != "", UserName: userName})
}

// HandleLeaderboard shows the dynamic restaurant list with scores.
func (s *Server) HandleLeaderboard(w http.ResponseWriter, r *http.Request) {
	restaurants, reviewsMap, _ := s.Repo.ListRestaurants(context.Background())

	config, _ := s.Repo.GetConfig(context.Background())
	weights := map[string]float64{"food": 0.5, "atmosphere": 0.2, "value": 0.2, "service": 0.1}
	if config != nil {
		weights = config.Weights
	}

	// Calculate scores for each restaurant
	for i := range restaurants {
		resReviews := reviewsMap[restaurants[i].ID]
		restaurants[i].OverallScore = scoring.CalculateRestaurantOverallScore(resReviews, weights)
		restaurants[i].ReviewCount = len(resReviews)
	}

	userID, userName := session.GetUser(r)
	s.render(w, "leaderboard.html", TemplateData{
		IsAuthenticated: userID != "",
		UserName:        userName,
		Restaurants:     restaurants,
		Weights:         weights,
	})
}

// HandleRestaurantDetail shows a single restaurant and its reviews.
func (s *Server) HandleRestaurantDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/restaurant/")
	id = strings.Split(id, "/")[0]

	restaurant, reviews, err := s.Repo.GetRestaurantData(context.Background(), id)
	if err != nil {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	config, _ := s.Repo.GetConfig(context.Background())
	weights := map[string]float64{"food": 0.5, "atmosphere": 0.2, "value": 0.2, "service": 0.1}
	if config != nil {
		weights = config.Weights
	}

	overall := scoring.CalculateRestaurantOverallScore(reviews, weights)
	userID, userName := session.GetUser(r)

	s.render(w, "restaurant_details.html", TemplateData{
		IsAuthenticated: userID != "",
		UserName:        userName,
		Restaurant:      restaurant,
		Reviews:         reviews,
		OverallScore:    overall,
	})
}

// HandleAddReview processes the review form.
func (s *Server) HandleAddReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, userName := session.GetUser(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id := parts[2]

	food, _ := strconv.ParseFloat(r.FormValue("food"), 64)
	atmosphere, _ := strconv.ParseFloat(r.FormValue("atmosphere"), 64)
	value, _ := strconv.ParseFloat(r.FormValue("value"), 64)
	service, _ := strconv.ParseFloat(r.FormValue("service"), 64)

	review := models.Review{
		PK:      fmt.Sprintf("RESTAURANT#%s", id),
		SK:      fmt.Sprintf("REVIEW#%s", userID),
		UserID:  userID,
		UserName: userName,
		Comment: r.FormValue("comment"),
		Date:    time.Now(),
		Ratings: map[string]float64{
			"food":       food,
			"atmosphere": atmosphere,
			"value":      value,
			"service":    service,
		},
	}

	err := s.Repo.SaveReview(context.Background(), review)
	if err != nil {
		http.Error(w, "Failed to save review", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/restaurant/"+id, http.StatusSeeOther)
}

// HandleLogin redirects to Google.
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state := "random-state"
	http.Redirect(w, r, s.AuthConfig.GetLoginURL(state), http.StatusTemporaryRedirect)
}

// HandleAuthCallback handles Google return.
func (s *Server) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	user, err := s.AuthConfig.VerifyUser(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	session.SetUser(w, r, user.ID, user.Name)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleLogout clears the user's session.
func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleAddRestaurant handles the form.
func (s *Server) HandleAddRestaurant(w http.ResponseWriter, r *http.Request) {
	userID, userName := session.GetUser(r)
	if userID == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		id := filepath.Base(filepath.Join("/", name))
		res := models.Restaurant{
			PK:           fmt.Sprintf("RESTAURANT#%s", id),
			SK:           "METADATA",
			ID:           id,
			Name:         name,
			Cuisine:      r.FormValue("cuisine"),
			Location:     r.FormValue("location"),
			GoogleMapURL: r.FormValue("maps_url"),
			ImageURL:     r.FormValue("image_url"),
		}
		s.Repo.SaveRestaurant(context.Background(), res)
		http.Redirect(w, r, "/leaderboard", http.StatusSeeOther)
		return
	}

	s.render(w, "add_restaurant.html", TemplateData{IsAuthenticated: true, UserName: userName})
}

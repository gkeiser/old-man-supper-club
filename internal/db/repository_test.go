package db

import (
	"testing"
	"time"

	"github.com/grant/old-man-supper-club/internal/models"
	"github.com/grant/old-man-supper-club/internal/scoring"
)

func TestScoringAndDataFetching(t *testing.T) {
	// 1. Setup Mock Weights
	weights := map[string]float64{
		"food":       0.5,
		"atmosphere": 0.2,
		"value":      0.2,
		"service":    0.1,
	}

	// 2. Setup Mock Reviews
	reviews := []models.Review{
		{
			Ratings: map[string]float64{
				"food":       10,
				"atmosphere": 8,
				"value":      8,
				"service":    7,
			},
			Comment: "Legendary mashed potatoes.",
			Date:    time.Now(),
		},
		{
			Ratings: map[string]float64{
				"food":       6,
				"atmosphere": 5,
				"value":      5,
				"service":    5,
			},
			Comment: "A bit disappointing today.",
			Date:    time.Now(),
		},
	}

	// 3. Test Weighted Scoring Logic
	// Review 1: (10*0.5) + (8*0.2) + (8*0.2) + (7*0.1) = 5.0 + 1.6 + 1.6 + 0.7 = 8.9
	// Review 2: (6*0.5) + (5*0.2) + (5*0.2) + (5*0.1) = 3.0 + 1.0 + 1.0 + 0.5 = 5.5
	// Overall Average: (8.9 + 5.5) / 2 = 7.2

	score1 := scoring.CalculateWeightedScore(reviews[0], weights)
	if score1 < 8.89 || score1 > 8.91 {
		t.Errorf("Expected score1 to be ~8.9, got %f", score1)
	}

	score2 := scoring.CalculateWeightedScore(reviews[1], weights)
	if score2 < 5.49 || score2 > 5.51 {
		t.Errorf("Expected score2 to be ~5.5, got %f", score2)
	}

	overall := scoring.CalculateRestaurantOverallScore(reviews, weights)
	if overall < 7.19 || overall > 7.21 {
		t.Errorf("Expected overall score to be ~7.2, got %f", overall)
	}
}

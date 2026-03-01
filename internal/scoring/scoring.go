package scoring

import "github.com/grant/old-man-supper-club/internal/models"

// CalculateWeightedScore computes the final score for a single review based on weights.
func CalculateWeightedScore(review models.Review, weights map[string]float64) float64 {
	var total float64
	for category, rating := range review.Ratings {
		if weight, ok := weights[category]; ok {
			total += rating * weight
		}
	}
	return total
}

// CalculateRestaurantOverallScore computes the average weighted score across all reviews.
func CalculateRestaurantOverallScore(reviews []models.Review, weights map[string]float64) float64 {
	if len(reviews) == 0 {
		return 0
	}

	var sum float64
	for _, review := range reviews {
		sum += CalculateWeightedScore(review, weights)
	}

	return sum / float64(len(reviews))
}

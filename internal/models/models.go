package models

import "time"

// Config represents global scoring weights
type Config struct {
	PK      string             `dynamodbav:"PK"` // CONFIG#GLOBAL
	SK      string             `dynamodbav:"SK"` // METADATA
	Weights map[string]float64 `dynamodbav:"Weights"`
}

// Restaurant represents a dining establishment
type Restaurant struct {
	PK           string `dynamodbav:"PK"` // RESTAURANT#<ID>
	SK           string `dynamodbav:"SK"` // METADATA
	ID           string `dynamodbav:"ID"`
	Name         string `dynamodbav:"Name"`
	Location     string `dynamodbav:"Location"`
	GoogleMapURL string `dynamodbav:"GoogleMapURL"`
	Cuisine      string  `dynamodbav:"Cuisine"`
	ImageURL     string  `dynamodbav:"ImageURL"`
	OverallScore float64 `dynamodbav:"-"` // Calculated field
	ReviewCount  int     `dynamodbav:"-"` // Calculated field
}

// Review represents a member's evaluation
type Review struct {
	PK       string             `dynamodbav:"PK"` // RESTAURANT#<ID>
	SK       string             `dynamodbav:"SK"` // REVIEW#<GOOGLE_ID>
	Ratings  map[string]float64 `dynamodbav:"Ratings"`
	Comment  string             `dynamodbav:"Comment"`
	Date     time.Time          `dynamodbav:"Date"`
	UserID   string             `dynamodbav:"UserID"`
	UserName string             `dynamodbav:"UserName"`
}

// User represents a club member
type User struct {
	PK        string `dynamodbav:"PK"` // USER#<GOOGLE_ID>
	SK        string `dynamodbav:"SK"` // METADATA
	GoogleID  string `dynamodbav:"GoogleID"`
	Email     string `dynamodbav:"Email"`
	Name      string `dynamodbav:"Name"`
	AvatarURL string `dynamodbav:"AvatarURL"`
	Role      string `dynamodbav:"Role"` // "admin" or "member"
}

// BlogPost represents a site update or article
type BlogPost struct {
	PK       string    `dynamodbav:"PK"` // BLOG#<ID>
	SK       string    `dynamodbav:"SK"` // METADATA
	Title    string    `dynamodbav:"Title"`
	Content  string    `dynamodbav:"Content"`
	AuthorID string    `dynamodbav:"AuthorID"`
	Date     time.Time `dynamodbav:"Date"`
}

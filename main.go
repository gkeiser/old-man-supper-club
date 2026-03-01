package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/grant/old-man-supper-club/internal/auth"
	"github.com/grant/old-man-supper-club/internal/db"
	"github.com/grant/old-man-supper-club/internal/session"
	"github.com/grant/old-man-supper-club/internal/web"
	"github.com/joho/godotenv"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed assets/*
var assetsFS embed.FS

func main() {
	// Load .env locally
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Session Store
	session.Init()

	// Initialize AWS Config & DynamoDB Client
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dbClient := dynamodb.NewFromConfig(cfg)
	repo := db.NewRepository(dbClient, os.Getenv("DYNAMO_TABLE_NAME"))

	// Initialize Auth Config
	authCfg := auth.NewConfig()
	server := web.NewServer(authCfg, repo, templateFS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Route definitions
	mux := http.NewServeMux()

	// Serve static files from embedded FS
	subAssets, _ := fs.Sub(assetsFS, "assets")
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(subAssets))))

	mux.HandleFunc("/", server.HandleHome)
	mux.HandleFunc("/leaderboard", server.HandleLeaderboard)
	mux.HandleFunc("/restaurant/", server.HandleRestaurantDetail)
	mux.HandleFunc("/restaurant/{id}/review", server.HandleAddReview)
	mux.HandleFunc("/login", server.HandleLogin)
	mux.HandleFunc("/auth/callback", server.HandleAuthCallback)
	mux.HandleFunc("/logout", server.HandleLogout)
	mux.HandleFunc("/add-restaurant", server.HandleAddRestaurant)

	// Middleware to force Content-Type for HTML routes
	htmlHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		mux.ServeHTTP(w, r)
	})

	// Check if running in Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Starting Lambda V2 adapter for Function URL")
		adapter := httpadapter.NewV2(htmlHandler)
		lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return adapter.ProxyWithContext(ctx, req)
		})
	} else {
		fmt.Printf("Old Man Supper Club server starting on http://localhost:%s\n", port)
		if err := http.ListenAndServe(":"+port, htmlHandler); err != nil {
			log.Fatal(err)
		}
	}
}

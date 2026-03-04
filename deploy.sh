#!/bin/bash

# Old Man Supper Club - Deploy Script
# Usage: ./deploy.sh

# 1. Exit on error
set -e

echo "🚀 Starting deployment for Old Man Supper Club..."

# 2. Compile for AWS Lambda (Linux/amd64)
echo "📦 Compiling Go binary..."
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go

# 3. Package into zip
echo "🗜️ Packaging into bootstrap.zip..."
zip -q bootstrap.zip bootstrap

# 4. Upload to AWS Lambda
echo "☁️ Uploading to AWS Lambda (us-east-2)..."
aws lambda update-function-code --function-name omsc-web-server --zip-file fileb://bootstrap.zip --region us-east-2

# 5. Clean up
echo "🧹 Cleaning up local build files..."
#rm bootstrap bootstrap.zip

echo "✅ Deployment Successful! Your changes are now live."

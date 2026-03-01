# Old Man Supper Club

A restaurant review site for a private club with weighted scoring.

## Project Guidelines
- **Backend:** Go (Standard Library preferred, but using AWS SDK & Gorilla Sessions).
- **Frontend:** Go `html/template` and Vanilla CSS.
- **Authentication:** Google OAuth2 with email whitelist lockdown.
- **Database:** DynamoDB (Single Table Design).
- **Deployment:** AWS Lambda + Function URLs (Preferred Region: us-east-2).
- **Build Mandate:** ALWAYS run `go mod tidy` and verify the project builds (`go build ./...`) before suggesting a run or deployment to the user.

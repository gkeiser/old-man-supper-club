# Old Man Supper Club - Technical Specification

## Project Overview
A restaurant review site for a private club, featuring weighted scoring, Google Authentication, and a $0-cost AWS serverless architecture.

## Authentication & Authorization
- [x] **Provider:** Google OAuth2.
- [x] **Lockdown Strategy:** 
    1. **Google Testing Mode:** Restricted to 100 specific emails in Google Console.
    2. **Application Whitelist:** Go middleware verifies email against `ALLOWED_EMAILS` env var.
- [x] **Session Management:** Secure, HTTP-only cookies via `gorilla/sessions`.

## Data Model (DynamoDB Single Table)
*All entities reside in a single table (e.g., 'SupperClub') using generic PK/SK attributes.*

### Partition Key (PK) | Sort Key (SK) | Attributes
| Entity | PK | SK | Attributes |
| :--- | :--- | :--- | :--- |
| **Global Config** | `CONFIG#GLOBAL` | `METADATA` | `Weights: {food: 0.5, value: 0.2, ...}` |
| **Restaurant** | `RESTAURANT#<ID>` | `METADATA` | `ID, Name, Address, GoogleMapURL, Cuisine, ImageURL` |
| **Review** | `RESTAURANT#<ID>` | `REVIEW#<GOOGLE_ID>` | `UserID, UserName, Ratings: {food: 9, ...}, Comment, Date` |
| **User** | `USER#<GOOGLE_ID>` | `METADATA` | `Email, Name, AvatarURL, Role` |
| **Blog Post** | `BLOG#<ID>` | `METADATA` | `Title, Content, AuthorID, Date` |

## Scoring Engine Logic
1. Fetch `CONFIG#GLOBAL` to get current category weights (Default: Food 0.5, Atmos 0.2, Value 0.2, Service 0.1).
2. Fetch all `REVIEW#` items for a `RESTAURANT#`.
3. For each review, calculate: `WeightedScore = Σ (CategoryRating * CategoryWeight)`.
4. Aggregate individual weighted scores for the restaurant's overall club rating.

## Roadmap & Status
- [x] **Phase 1: Project Foundation** (Go modules, AWS SDK, structure).
- [x] **Phase 2: Core Data Layer** (DynamoDB Repo, Single Table Design).
- [x] **Phase 3: Authentication** (Google OAuth2 + Whitelist + Sessions).
- [x] **Phase 4: Restaurant & Review Features** (Add Restaurant, Detail View, Submit Review).
- [x] **Phase 5: Scoring Engine** (Weighted average calculation logic).
- [ ] **Phase 6: Blog System** (CRUD for club updates/articles).
- [ ] **Phase 7: AWS Deployment** (SAM/CDK for Lambda + Function URLs + CloudFront).
- [ ] **Phase 8: UI Polish** (Final styling, mobile responsiveness).

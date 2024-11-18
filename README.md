# Gocial
Gocial is an example social media application that abides by the 12 factor app conventions. Its a simple yet complete demonstration to be forked, copied, and extended.

### Features
#### Completed
- Postgresql Transactions
- Migrations
- CRUD Operations
- User feeds
- Pagination, filtering and sorting
- Generated Swagger documentation
- Structured Logging
- Email confirmation and verification

#### In Progress
- Authentication (JWT)
- Authorization
- Redis Caching
- Testing
- Graceful Shutdowns
- Rate Limiting
- Server Metrics
- Automation (CI/CD)




### Running Migrations

`migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_users`

`migrate -path=./cmd/migrate/migrations -database="postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable" up`
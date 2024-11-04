### Running Migrations

`migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_users`

`migrate -path=./cmd/migrate/migrations -database="postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable" up`
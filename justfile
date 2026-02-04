migrate-up:
    migrate -database sqlite3://app.db -path migrations up

migrate-down:
    migrate -database sqlite3://app.db -path migrations down

test:
    go test -v ./...

test-coverage:
    go test -cover ./...

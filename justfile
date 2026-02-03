migrate-up:
    migrate -database sqlite3://app.db -path migrations up

migrate-down:
    migrate -database sqlite3://app.db -path migrations down

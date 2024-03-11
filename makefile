setup-db:
	sqlite3 messages.db < init_db.sql

setup-test-db: init_db.sql
	sqlite3 test.db < init_db.sql

run: format vet
	go run main.go

test: format vet setup-test-db
	go test -v ./...

vet:
	go vet ./...

format:
	gofmt -w .
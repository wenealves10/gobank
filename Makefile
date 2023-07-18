createmigration:
	migrate create -ext=sql -dir=./sql/migrations -seq $(word 2,$(MAKECMDGOALS))
	@echo "Migration created: $(word 2,$(MAKECMDGOALS))"
%:
	@:

migrate:
	migrate -path=./sql/migrations -database "YOUR_URL_DATABASE_POSTGRES" -verbose up

migratedown:
	migrate -path=./sql/migrations -database "YOUR_URL_DATABASE_POSTGRES" -verbose down

dev:
	go run main.go

build:
	go build .

debug:
	go run -gcflags "-N -l"  cmd/chatservice/main.go

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createmigration migrate migratedown dev build debug sqlc test
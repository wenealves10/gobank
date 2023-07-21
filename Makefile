createmigration:
	migrate create -ext=sql -dir=./sql/migrations -seq $(word 2,$(MAKECMDGOALS))
	@echo "Migration created: $(word 2,$(MAKECMDGOALS))"
%:
	@:

migrateup:
	migrate -path=./sql/migrations -database "postgres://gobank:gobank1234@localhost:5434?sslmode=disable&database=gobank" -verbose up

migratedown:
	migrate -path=./sql/migrations -database "postgres://gobank:gobank1234@localhost:5434?sslmode=disable&database=gobank" -verbose down

dev:
	go run main.go

build:
	go build .

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createmigration migrateup migratedown dev build sqlc test
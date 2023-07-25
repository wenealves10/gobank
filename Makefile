createmigration:
	migrate create -ext=sql -dir=./sql/migrations -seq $(word 2,$(MAKECMDGOALS))
	@echo "Migration created: $(word 2,$(MAKECMDGOALS))"
%:
	@:

migrateup:
	migrate -path=./sql/migrations -database "$(DB_SOURCE)" -verbose up

migratedown:
	migrate -path=./sql/migrations -database "$(DB_SOURCE)" -verbose down

dev:
	go run main.go

build:
	go build .

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -source=db/sqlc/store.go -package=mocks -destination=db/mocks/store_mock.go

.PHONY: createmigration migrateup migratedown dev build sqlc test server mock
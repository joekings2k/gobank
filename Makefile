-include app.env
export $(shell sed 's/=.*//' .env)

postgres:
	docker run --name postgres12 --network bank-network -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p 5432:5432 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=$(DB_USER) --owner=root simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_LOCAL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_LOCAL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_LOCAL)" -verbose down 1

dropdb:
	docker exec -it postgres12 dropdb simple_bank

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run $(MAIN_GO)

mock:
	mockgen -package=mockdb -destination=db/mock/store.go --build_flags=--mod=mod github.com/joekings2k/gobank/db/sqlc Store

.PHONY:postgres  createdb migrateup migratedown migrateup1 migratedown1 dropdb sqlc test server mock



# a69e9134296d24094b9dba6bc466c90f-716640240.us-east-1.elb.amazonaws.com
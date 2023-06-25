
pwd = C:\Users\morka_joshua\StudioProjects\GoProjects\shinybank

postgres: 
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine


createdb:
	docker exec -it postgres12 createdb --username=root --owner=root shiny_bank

dropdb:
	docker exec -it postgres12 dropdb shiny_bank

createmigrate:
	migrate create -ext sql -dir db/migration -seq init_schema


migrateup:
	migrate -path  ./src/db/migration -database "postgresql://root:secret@localhost:5432/shiny_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path  ./src/db/migration -database "postgresql://root:secret@localhost:5432/shiny_bank?sslmode=disable" -verbose up 1


migratedown:
	migrate -path  src/db/migration -database "postgresql://root:secret@localhost:5432/shiny_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path  src/db/migration -database "postgresql://root:secret@localhost:5432/shiny_bank?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate
	# sqlc generate

mockgenstore:
	mockgen -package mockdb -destination src/db/mock/store.go  github.com/morka17/shiny_bank/v1/src/db/sqlc Store

server:
	go run main.go

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown1 migratedown sqlc test server
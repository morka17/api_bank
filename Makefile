
pwd = C:\Users\morka_joshua\StudioProjects\GoProjects\shinybank
# dburl = shine-bank.cwaywcycutdv.us-east-1.rds.amazonaws.com
dburl = localhost:5432
dbname=shiny_bank

build:
	CGO_ENABLED=0 GOOS=linux go build -o main

run: build 
	@./bin/microservices2

postgres: 
	docker run --name postgres12 --network api-bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root shiny_bank

dropdb:
	docker exec -it postgres12 dropdb shiny_bank

createmigrate:
	migrate create -ext sql -dir ./src/db/migration -seq $(name)

migrateup:
	migrate -path  ./src/db/migration -database "postgresql://root:secret@$(dburl)/$(dbname)?sslmode=disable" -verbose up

migrateup1:
	migrate -path  ./src/db/migration -database "postgresql://root:secret@$(dburl)/$(dbname)?sslmode=disable" -verbose up 1

migrateup2:
	migrate -path  ./src/db/migration -database "postgresql://root:secret@$(dburl)/$(dbname)?sslmode=disable" -verbose up 2

migratedown:
	migrate -path  src/db/migration -database "postgresql://root:secret@$(dburl)/$(dbname)?sslmode=disable" -verbose down

migratedown1:
	migrate -path  src/db/migration -database "postgresql://root:secret@$(dburl)/$(dbname)?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate
	# sqlc generate

mockgenstore:
	mockgen -package mockdb -destination src/db/mock/store.go  github.com/morka17/shiny_bank/v1/src/db/sqlc Store
	mockgen -package mockwk -destination src/worker/mock/distributor.go  github.com/morka17/shiny_bank/v1/src/worker TaskDistributor 

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

server:
	go run main.go

docker:
	docker build -t shinybank:latest .	

dockerContainer:
	docker run --name shinybank --network api-bank-network -p 8080:8080 -e GIN_MODE=release -e DB_Source="postgresql://root:secret@postgres12:5432/shiny_bank?sslmode=disable" shinybank:latest 

createdockernetwork:
	docker network create api-bank-network

connectdockernetwork:
	docker network connect api-bank-network postgres12

proto:	
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	# protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	# --go-grpc_out=pb  --go-grpc_opt=paths=source_relative  \
	# --go-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	# proto/*.proto
	protoc -I ./proto \
  	--go_out ./pb --go_opt paths=source_relative \
  	--go-grpc_out ./pb --go-grpc_opt paths=source_relative \
  	--grpc-gateway_out ./pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out ./doc/swagger --openapiv2_opt allow_merge=true,merge_file_name=shiny_bank\
  	./proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans: 
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379  -d  redis:7-alpine  


test:
	go test -v -cover -short ./...

.PHONY: postgres mockgenstore createdb db_docs db_schema dropdb migrateup migrateup1 migratedown1 migratedown sqlc test server proto
build:
	go build -o ./backend .

fake:
	go build -o ./fake_client github.com/FocusCompany/backend-go/fake


DB_ADDR ?= 127.0.0.1
DB_PORT ?= 5432
migrate:
	migrate -source file://database/ -database "postgres://postgres@?sslmode=disable&host=$(DB_ADDR)&port=$(DB_PORT)" up
migrate-down:
	migrate -source file://database/ -database "postgres://postgres@?sslmode=disable&host=$(DB_ADDR)&port=$(DB_PORT)" down
drop-db:
	migrate -source file://database/ -database "postgres://postgres@?sslmode=disable&host=$(DB_ADDR)&port=$(DB_PORT)" drop
migrate-new:
	migrate create -dir database -ext sql $(FILE)

proto:
	protoc --proto_path=./protobuf_envelope --go_out=:./proto protobuf_envelope/*.proto


.DEFAULT_GOAL=build

.PHONY: fake proto
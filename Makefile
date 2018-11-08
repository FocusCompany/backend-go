build:
	go build -o ./backend .

fake:
	go build -o ./fake_client github.com/FocusCompany/backend-go/fake

migrate:
	migrate -source file://database/ -database "postgres://postgres@127.0.0.1:5432?sslmode=disable" up

migrate-down:
	migrate -source file://database/ -database "postgres://postgres@127.0.0.1:5432?sslmode=disable" down

drop-db:
	migrate -source file://database/ -database "postgres://postgres@127.0.0.1:5432?sslmode=disable" drop

migrate-new:
	migrate create -dir database -ext sql $(FILE)

proto:
	protoc --proto_path=./protobuf_envelope --go_out=:./proto protobuf_envelope/*.proto


.DEFAULT_GOAL=build

.PHONY: fake proto
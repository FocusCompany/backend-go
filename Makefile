build:
	go build -o ./backend .

migrate:
	migrate -source file://database/ -database "postgres://postgres@127.0.0.1:5432?sslmode=disable" up

migrate-down:
	migrate -source file://database/ -database "postgres://postgres@127.0.0.1:5432?sslmode=disable" down

drop-db:
	migrate -source file://database/ -database "postgres://postgres@127.0.0.1:5432?sslmode=disable" drop

migrate-new:
	migrate create -dir database -ext sql $(FILE)

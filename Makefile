DB_URL=postgresql://root:secret@localhost:5433/chain_task?sslmode=disable

run: 
	go run main.go

docker:
	docker build -t chaintask:latest .

network:
	docker network create task-network

postgres:
	docker run --name taskchain-postgres --network task-network -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine


createdb:
	docker exec -it taskchain-postgres createdb --username=root --owner=root chain_task

dropdb:
	docker exec -it taskchain-postgres dropdb chain_task

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

test:
	go test -v -cover -short ./...


mock:
	mockgen -package mockdb --build_flags=--mod=mod -destination internal/mock/taskStore.go github.com/kamalbowselvam/chaintask/internal/core/ports TaskRepository
	mockgen -package mockdb --build_flags=--mod=mod -destination internal/mock/userStore.go github.com/kamalbowselvam/chaintask/internal/core/ports UserRepository

.PHONY: run network postgres createdb dropdb migrateup migratedown test mock
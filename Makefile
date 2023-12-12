DB_URL=postgresql://root:secret@localhost:5432/chain_task?sslmode=disable

server: 
	go run main.go

docker:
	docker build -t chaintask:latest .

network:
	docker network create task-network

postgres:
	docker run --name taskchain-postgres --network task-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

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

swagger:
	swag init

mock:
	mockgen -package mockdb --build_flags=--mod=mod -destination mock/store.go github.com/kamalbowselvam/chaintask/db GlobalRepository

kind: cluster docker upload deploy 

cluster:
	kind create cluster --config=./kind/clusterconfig.yaml && chmod 777 ./kind/registry.sh && ./kind/registry.sh

upload: 
	docker tag chaintask:latest localhost:5001/chaintask:latest && 	docker push localhost:5001/chaintask:latest

deploy:
	kubectl apply -f ./kind/postgres.yaml && sleep 10 && kubectl apply -f ./kind/chaintask.yaml

undeploy:
	kind delete cluster 

.PHONY: server network postgres createdb dropdb migrateup migratedown test mock migrateawsup migrateawsdown swagger cluster upload deploy undeploy

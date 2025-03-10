run:
	go run cmd/http/main.go

build:
	@go build -o bin/go-parking-lot-users-service cmd/http/main.go

docker_build:
	docker build -t go-parking-lot-users-service . 

test:
	@go test -v ./...

sqlc:
	sqlc generate

migrateup:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/postgres/migration -seq $(name)


k8s_apply: 
	kubectl apply -f ../deployments/kubernetes/users/users-service-deployment.yaml

k8s_rollout:
	kubectl rollout restart deployment/go-parking-lot-users-service

create_configmap:
	kubectl create configmap users-config --from-env-file=../deployments/kubernetes/users/.env

.PHONY:
	run build test sqlc migrateup migrateup1 migratedown migratedown1 new_migration 

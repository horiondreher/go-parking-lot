run:
	go run cmd/http/main.go

build:
	@go build -o bin/go-parking-lot-parkings-service cmd/http/main.go

test:
	@go test -v ./...

create_configmap:
	kubectl create configmap parkings-config --from-env-file=../deployments/kubernetes/parkings/.env

.PHONY:
	run build test sqlc migrateup migrateup1 migratedown migratedown1 new_migration 

run:
	go run cmd/http/main.go

build:
	@go build -o bin/go-parking-lot-parkings-service cmd/http/main.go

docker_build:
	docker build -t go-parking-lot-parkings-service . 

test:
	@go test -v ./...

k8s_apply: 
	kubectl apply -f ../deployments/kubernetes/parkings/parkings-service-deployment.yaml 

k8s_rollout:
	kubectl rollout restart deployment/go-parking-lot-parkings-service

create_configmap:
	kubectl create configmap parkings-config --from-env-file=../deployments/kubernetes/parkings/.env

.PHONY:
	run build test sqlc migrateup migrateup1 migratedown migratedown1 new_migration 

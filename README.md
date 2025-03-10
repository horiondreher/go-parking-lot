# Go Parking Lot Design

## Deploy services in Kubernetes

Store a `.env` in both services `/deployments/kubernetes/users/.env` and `/deployments/kubernetes/parkings/.env`:

### Users Service

```sh
ENVIRONMENT=development

HTTP_SERVER_ADDRESS=0.0.0.0:8080
QUEUE_SERVER_ADDRESS=amqp://guest:guest@rabbitmq:5672/

POSTGRES_DB=go_parking_lot_user_service
POSTGRES_USER=pguser
POSTGRES_PASSWORD=pgpassword

MIGRATION_URL=

TOKEN_SYMMETRIC_KEY=
ACCESS_TOKEN_DURATION=
REFRESH_TOKEN_DURATION=
```

### Parkings Service

```sh
HTTP_SERVER_ADDRESS=0.0.0.0:8080
GRPC_SERVER_ADDRESS=0.0.0.0:50051
QUEUE_SERVER_ADDRESS=amqp://guest:guest@rabbitmq:5672/
```

Start `minikube`

```sh
minikube start
```

Create the RabbitMQ service

```sh
kubectl apply -f deployments/kubernetes/rabbitmq/rabbitmq-deployment.yaml
```

Build the service image and set variables for `minikube` in each folder (`/users` and `/parkings`):

```sh
eval $(minikube docker-env)
docker build -t go-parking-lot-users-service .
docker build -t go-parking-lot-parkings-service .
```

Apply the k8s specs files for users-service

```sh
make create_configmap
kubectl apply -f deployments/kubernetes/users/persistent-volume.yaml
kubectl apply -f deployments/kubernetes/users/postgres-deployment.yaml
kubectl apply -f deployments/kubernetes/users/users-service-deployment.yaml
```

Apply the k8s specs files for parkings-service

```sh
make create_configmap
kubectl apply -f deployments/kubernetes/parkings/parkings-service-deployment.yaml
```

Deploy the ingress to access users-service externally

```sh
kubectl apply -f deployments/kubernetes/parkings-log-ingress.yaml
```

## Other commands

### Delete configmaps

```sh
kubectl delete configmaps
```

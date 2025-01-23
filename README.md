
# Go Parking Lot Design

## Deploy services in Kubernetes

Store a `.env` in `deployments/kubernetes/.env`:

```sh
ENVIRONMENT=

HTTP_SERVER_ADDRESS=

POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=

MIGRATION_URL=

TOKEN_SYMMETRIC_KEY=
ACCESS_TOKEN_DURATION=
REFRESH_TOKEN_DURATION=
```

Start `minikube`

```sh
minikube start
```

Build the service image and set variables for `minikube`:

```sh
docker build -t go-parking-lot-user-service .
eval $(minikube docker-env)
```

Apply the k8s specs files

```sh
make create_configmap
kubectl apply -f deployments/kubernetes/persistent-volume.yaml
kubectl apply -f deployments/kubernetes/postgres-deployment.yaml
kubectl apply -f deployments/kubernetes/user-service-deployment.yaml
```

## Other commands

### Delete configmaps

```sh
kubectl delete configmaps
```

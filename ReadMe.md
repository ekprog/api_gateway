protoc -I ./proto/auth_service \
--go_out ./pkg/auth_service \
--go_opt paths=source_relative \
--go-grpc_out ./pkg/auth_service \
--go-grpc_opt paths=source_relative \
--grpc-gateway_out ./pkg/auth_service \
--grpc-gateway_opt paths=source_relative \
./proto/auth_service/api/*.proto



## 1. Build docker
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./server
docker build -t api_gw3 .
```


## 2. Push to docker hub
```bash
#docker login -u <username>
docker tag api_gw3 egorkozelskij/api_gw3
docker push egorkozelskij/api_gw3
```


## For development run
```bash
docker compose --env-file .env up postgres
docker compose --env-file .env up api_gw
```

## 3. Deploy
```bash
docker compose --env-file .env up
```
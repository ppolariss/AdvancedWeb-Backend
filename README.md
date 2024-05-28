# AdvancedWeb-Backend
Advanced Web Backend

## build and run
1. 
```
cd backend
go run main.go

cd driving
go run .
```

2. 
```
cd backend
go install github.com/swaggo/swag/cmd/swag@latest
swag init --parseDependency --parseDepth 1
go build -o main.exe
./main.exe
```

3. 
```
<!-- docker pull ppolariss/advanced_web_backend:latest -->
docker compose -p awb pull
docker compose -p awb up -d
```

GOOS=linux go build
docker build -t rjames187/gateway:1.0 .
go clean
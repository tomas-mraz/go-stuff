
set GOOS=linux
set GOARCH=amd64

go mod tidy

go build -o server grpc2/server

go build -o client grpc2/client

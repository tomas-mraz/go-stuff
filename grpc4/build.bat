
set GOOS=linux
set GOARCH=amd64

go mod tidy

go build -o server grpc3/server

go build -o client grpc3/client

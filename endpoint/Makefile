
all: windows linux

windows:
	GOOS=windows \
	GOARCH=amd64 \
	CGO_ENABLED=1 \
	go build -ldflags="-H windowsgui"

linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build

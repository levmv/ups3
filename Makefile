
build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ups3-w64 main.go
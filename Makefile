build:
	go build -o "goors-cli"

build-windows:
	CGO_ENABLED=0 \
	GOOS=windows \
	GOARCH=amd64 \
	go build \
	-o "goors-cli-windows-x64.exe" \
	main.go

build-linux:
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build \
	-o "goors-cli-linux-x64" \
	main.go

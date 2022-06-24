BOT_BINARY = realty_bot

build_linux:
	@echo "Building linux bot..."
	env GOOS=linux CGO_ENABLED=0 go build -ldflags "-s -w" -v -o build/${BOT_BINARY}_linux ./cmd/app
	@echo "Done!"

build_macos:
	@echo "Building MacOS bot..."
	env GOOS=darwin CGO_ENABLED=0 go build -ldflags "-s -w" -v -o build/${BOT_BINARY}_macos ./cmd/app
	@echo "Done!"

build_pi:
	@echo "Building MacOS bot..."
	env GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -ldflags "-s -w" -v -o build/${BOT_BINARY}_pi ./cmd/app
	@echo "Done!"

test:
	@echo "Run tests..."
	go test -v ./internal/...

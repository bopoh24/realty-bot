BOT_BINARY = realty_bot

build.linux:
	@echo "Building linux bot..."
	env GOOS=linux CGO_ENABLED=0 go build -ldflags "-s -w" -v -o build/${BOT_BINARY}_linux ./cmd/app
	@echo "Done!"

build.macos:
	@echo "Building MacOS bot..."
	env GOOS=darwin CGO_ENABLED=0 go build -ldflags "-s -w" -v -o build/${BOT_BINARY}_macos ./cmd/app
	@echo "Done!"

build.pi:
	@echo "Building MacOS bot..."
	env GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -ldflags "-s -w" -v -o build/${BOT_BINARY}_pi ./cmd/app
	@echo "Done!"

test:
	@echo "Run tests..."
	go test -v ./internal/...

up:
	@echo "Starting app in docker..."
	docker-compose up -d
	@echo "Dockerized app started!"

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

up.build:
	@echo "Rebuild and starting app in docker..."
	docker-compose up --build -d
	@echo "Dockerized app started!"


up.pi:
	@echo "Starting app in docker..."
	docker-compose -f docker-compose-pi.yml up -d
	@echo "Dockerized app started!"

down.pi:
	@echo "Stopping docker compose..."
	docker-compose -f docker-compose-pi.yml down
	@echo "Done!"

up.build.pi:
	@echo "Rebuild and starting app in docker..."
	docker-compose -f docker-compose-pi.yml up --build -d
	@echo "Dockerized app started!"
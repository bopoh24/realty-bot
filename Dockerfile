##
## Build
##
FROM golang:1.18.3-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -ldflags "-s -w" -v -o realty_bot ./cmd/app

##
## Deploy
##
FROM alpine:latest

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/realty_bot /app/realty_bot

WORKDIR /app

CMD ["/app/realty_bot"]

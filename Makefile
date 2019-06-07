install:
	# Go client
	@go get github.com/nats-io/nats.go/

	# Server
	@go get github.com/nats-io/nats-server

	# Streaming
	@go get github.com/nats-io/stan.go/

up:
	@docker-compose up -d

down:
	@docker-compose down

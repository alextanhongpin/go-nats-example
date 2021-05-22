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

util:
	@docker run --network nats --rm -it synadia/nats-box:latest
	@# For every command in this shell, specify the server name, e.g.
	@# nats str ls -s nats://nats-main:4222
	@# Dropping the port is okay.
	@# nats str ls -s nats://nats-main

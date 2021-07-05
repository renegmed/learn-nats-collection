init-project:
	go mod init practical-nats/driver-agent
 

build:
	go build -race -o driver-agent cmd/main.go

run: build 
	./driver-agent -s nats://localhost:4222

build-docker:
	docker build . --tag renegmedal/trans-service-agent:v0.0.1 

push:
	docker push renegmedal/trans-service-agent:v0.0.1 

# run-server:

# 	go run cmd/main.go --nats "nats://localhost:4221,nats://localhost:4222,nats://localhost:4223" --listen localhost:8090 
# .PHONY: run-server

# server-client:
# 	curl localhost:8090/
# .PHONY: server-client

# rider-request:
# 	curl -X POST localhost:8080/rides
# .PHONY: rider-request
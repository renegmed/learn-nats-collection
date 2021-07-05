init-project:
	go mod init practical-nats/riders-client

build:
	go build -race -o api-server cmd/main.go

run: build 
	./api-server -s nats://localhost:4222 -port 8282 

build-docker:
	docker build . --tag renegmedal/trans-service-api:v0.0.1 

push:
	docker push renegmedal/trans-service-api:v0.0.1 
 
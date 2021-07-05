 init-project:
	go mod init practical-nats/rides-manager
 
build:
	go build -race -o rides-manager cmd/main.go

run: build 
	./rides-manager -s nats://localhost:4222

build-docker:
	docker build . --tag renegmedal/trans-service-manager:v0.0.1 

push:
	docker push renegmedal/trans-service-manager:v0.0.1 
init-project:
	go mod init jetstream-order 
 
build:
	go build -race -o new-order .

run: build
	./new-order -s nats://localhost:4222


build-docker:
	docker build . --tag renegmedal/jetstream-order:v0.0.1

push:
	docker push renegmedal/jetstream-order:v0.0.1
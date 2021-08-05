init-project:
	go mod init jetstream-monitor 

build:
	go build -race -o monitor .

run: build
	./monitor -s nats://localhost:4222

build-docker:
	docker build . --tag renegmedal/jetstream-monitor:v0.0.1

push:
	docker push renegmedal/jetstream-monitor:v0.0.1
	
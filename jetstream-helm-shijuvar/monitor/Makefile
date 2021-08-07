init-project:
	go mod init ordering-app 

build:
	go build -race -o monitor .

run: build
	./monitor -s nats://localhost:4222

build-docker:
	docker build . --tag renegmedal/ordering-monitor:v0.0.1

push:
	docker push renegmedal/ordering-monitor:v0.0.1
	
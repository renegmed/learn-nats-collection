init-project:
	go mod init jetstream-review 

build:
	go build -race -o order-review .

run: build
	./order-review -s nats://localhost:4222
	

build-docker:
	docker build . --tag renegmedal/jetstream-review:v0.0.1

push:
	docker push renegmedal/jetstream-review:v0.0.1
	
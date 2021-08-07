init-project:
	go mod init ordering-app 
 
build:
	go build -race -o new-order .

run: build
	./new-order -s nats://localhost:4222


build-docker:
	docker build . --tag renegmedal/ordering-new:v0.0.2

push:
	docker push renegmedal/ordering-new:v0.0.2
	
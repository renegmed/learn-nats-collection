init-project:
	go mod init ordering-app 

build:
	go build -race -o order-review .

run: build
	./order-review -s nats://localhost:4222
	

build-docker:
	docker build . --tag renegmedal/ordering-review:v0.0.2

push:
	docker push renegmedal/ordering-review:v0.0.2
	
	
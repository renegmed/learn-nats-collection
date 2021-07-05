init-project:
	go mod init product-read 

build:
	go build -race -o read . 

run: build 
	./read -port 8081 -s nats://localhost:4222 -durablename cqrs-read -topic PRODUCT.inserted 

build-docker:
	docker build . --tag renegmedal/cqrs-product-read:v0.0.1

push:
	docker push renegmedal/cqrs-product-read:v0.0.1

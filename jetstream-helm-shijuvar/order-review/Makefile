init-project:
	go mod init jsdemo-shijuvar 

build:
	go build -race -o order-review .

run: build
	./order-review -s nats://localhost:4222
	

	
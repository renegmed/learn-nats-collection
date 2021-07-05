init-project:
	go mod init jsdemo-shijuvar 

build:
	go build -race -o monitor .

run: build
	./monitor -s nats://localhost:4222


	
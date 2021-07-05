init-project:
	go mod init jsdemo-shijuvar 
 
build:
	go build -race -o new-order .

run: build
	./new-order -s nats://localhost:4222
	
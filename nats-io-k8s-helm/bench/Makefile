init-project:
	go mod init nats-bench 

build: 
	go build -race -o bench . 

run: build 
	./bench -s nats://localhost:4222
init-project:
	go mod init nats-echo 

build:
	go build -race -o echo . 

run: build 
	./echo -s nats://localhost:4222 test 

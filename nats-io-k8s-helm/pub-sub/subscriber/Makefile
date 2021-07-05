init-project:
	go mod init nats-subscriber 

build:
	go build -race -o subscriber .

run: build
	./subscriber -s nats://localhost:4222 test 


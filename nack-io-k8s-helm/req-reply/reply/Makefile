init-project:
	go mod init nats-reply 

build:
	go build -race -o reply . 

run: build 
	./reply -s nats://localhost:4222 -q products itemA.inventory.status  105

run2: build 
	./reply -s nats://localhost:4222 -q products itemA.inventory.status  260
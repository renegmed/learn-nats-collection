init-project:
	go mod init nats-reply 

build:
	go build -race -o reply . 

run: build 
	./reply -s nats://localhost:4222 itemA.inventory.status  105


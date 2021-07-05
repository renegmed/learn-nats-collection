init-project:
	go mod init nats-request 

build:
	go build -race -o request .

run: build 
	./request -s nats://localhost:4222 itemA.inventory.status 'Provide current inventory of item A.'
init-project:
	go mod init nats-qpub 

build:
	go build -race -o qpub . 

run: build 
	./qpub -s nats://127.0.0.1:4222 Topic.A 'This is a queue test.'
 
run2: build 
	./qpub -s nats://127.0.0.1:4222 Topic.A 'THIS IS ANOTHER QUEUE TEST'
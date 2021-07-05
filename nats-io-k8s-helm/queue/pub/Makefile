init-project:
	go mod init nats-qpub 

build:
	go build -race -o qpub . 

run: build 
	./qpub -s nats://127.0.0.1:4222 test 'This is a queue test.'


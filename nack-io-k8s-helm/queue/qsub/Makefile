init-project:
	go mod init nats-qsub 

build:
	go build -race -o qsub .

run: build 
	./qsub -s nats://localhost:4222 Topic.A topic-queue

run2: build 
	./qsub -s nats://localhost:4222 Topic.A topic-queue
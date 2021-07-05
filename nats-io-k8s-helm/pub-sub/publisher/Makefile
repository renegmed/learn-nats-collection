init-project:
	go mod init nats-publisher

build:
	go build -race -o publisher .


run: build
	./publisher -s nats://localhost:4222 test 'hello there!'


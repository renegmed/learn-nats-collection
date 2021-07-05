init-project:
	go mod init nats-publisher

build:
	go build -race -o publisher .


run: build
	./publisher -s nats://localhost:4222 TOPIC.a 'This is the message!'

run2: build
	./publisher -s nats://localhost:4222 TOPIC.a 'THIS IS A NEW MESSAGE!'

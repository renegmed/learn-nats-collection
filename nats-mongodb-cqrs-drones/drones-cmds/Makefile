init-project:
	go mod init cqrs-drones/cmds

build:
	go build -race -o cmds . 

run: build 
	./cmds -s nats://localhost:4222 -port 7070


build-docker:
	# rm cmds 
	docker build . --tag renegmedal/drones-commands:v0.0.1

push: 
	docker push renegmedal/drones-commands:v0.0.1
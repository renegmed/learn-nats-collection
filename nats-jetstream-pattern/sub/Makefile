init-project:
	go mod init nats-stream-pattern
 
build:
	go build -race -o subscriber . 

run: build  
	./subscriber -s nats://localhost:4222 -topic ORDER.created -durablename dur-sub

run2: build  
	./subscriber -s nats://localhost:4222 -topic ORDER.update -durablename dur-change


# build and push image to docker hub
build-docker:
	docker build -t renegmedal/js-pattern-subscriber:0.1.0 .  

push:
	docker push renegmedal/js-pattern-subscriber:0.1.0 

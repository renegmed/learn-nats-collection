#export NATS_URL="nats://172.22.0.3:4222" #,nats://172.22.0.2:4222,nats://172.22.0.4:4222"
#export NATS_ADDR=nats://nats1:4222 #,nats://172.22.0.2:4222,nats://172.22.0.4:4222"
export NATS_ADDR=nats://localhost:4222
export NATS_CLUSTER_ID=cluster-one
export NATS_CLIENT_ID=pub-1
export SERVER_PORT=9090

init-project:
	go mod init nats-stream-pattern

build:
	go build -race -o publisher .

run: build 
	./publisher -s nats://localhost:4222 -streamname ORDER -topic ORDER.* -port 9090 


post1:
	curl -X POST http://localhost:9090/publish -d '{"topic":"ORDER.created", "message":"Order 10013"}'


post2:
	curl -X POST http://localhost:9090/publish -d '{"topic":"ORDER2.cancelled", "message":"Order 1015"}'

 

# build and push image to docker hub
build-docker:
	docker build -t renegmedal/js-pattern-publisher:0.1.0 .  

push:
	docker push renegmedal/js-pattern-publisher:0.1.0 


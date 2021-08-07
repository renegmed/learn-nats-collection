init-project:
	go mod init cqrs-drones/query

build:
	go build -race -o query .

run: build 
	./query -s nats://localhost:4222 


build-docker:
	# rm query 
	docker build . --tag renegmedal/drones-query:v0.0.1

push: 
	docker push renegmedal/drones-query:v0.0.1
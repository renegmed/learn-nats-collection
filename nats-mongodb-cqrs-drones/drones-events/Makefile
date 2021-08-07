init-project:
	go mod init cqrs-drones/events

build: 
	go build -race -o events . 

run: build 
	./events -s nats://localhost:4222 -dburl mongodb://localhost:27017
	

build-docker:
	# rm events 
	docker build . --tag renegmedal/drones-events:v0.0.1

push: 
	docker push renegmedal/drones-events:v0.0.1
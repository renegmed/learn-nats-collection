init-project:
	go mod init cinema-app/users 

build:
	go build -race -o users cmd/app/* 

run: build
	./users -serveraddr localhost:7000 -mongouri mongodb://localhost:27020 -enablecredentials=false -s nats://localhost:4222 -t=USERS -q qusers 


build-docker:
	docker build . --tag renegmedal/cinema-nats-users:latest

push:
	docker push renegmedal/cinema-nats-users:latest


get-all:
	curl localhost:7000/api/users/ -v

post1:
	curl -XPOST localhost:7000/api/users/ -d '{"name":"John","lastname":"Doe"}' -v 

post2:
	curl -XPOST localhost:7000/api/users/ -d '{"name":"Jane","lastname":"Denali"}'  
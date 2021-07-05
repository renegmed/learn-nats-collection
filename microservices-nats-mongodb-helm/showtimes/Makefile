init-project:
	go mod init cinema-app/showtimes

build:
	go build -race -o showtimes cmd/app/* 

run: build
	./showtimes -serveraddr localhost:6000 -mongouri mongodb://localhost:27019 -enablecredentials=false -s nats://localhost:4222 -t=SHOWTIMES -q qshowtimes
 

up:
	docker-compose up -d --build

run-docker:
	curl localhost:8080/

build-docker:
	docker build . --tag renegmedal/cinema-nats-showtimes:latest

push:
	docker push renegmedal/cinema-nats-showtimes:latest


get-all:
	curl localhost:6000/api/showtimes/ -v

post1:
	curl -XPOST localhost:6000/api/showtimes/ -d '{"date":"2021-07-01", "movies":["60e272e5ac81947e79847c17","60e2bff5ac81947e79847c19"]}' -v



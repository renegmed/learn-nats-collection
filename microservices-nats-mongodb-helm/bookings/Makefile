init-project:
	go mod init cinema-app/bookings

build:
	go build -race -o bookings cmd/app/* 

run: build
	./bookings -serveraddr localhost:4000 -mongouri mongodb://localhost:27017 -enablecredentials=false -s nats://localhost:4222 -t=BOOKINGS -q qbookings


get-all:
	curl http://localhost:4000/api/bookings/ -v

post1:
	curl -XPOST http://localhost:4000/api/bookings/ -d '{"userid":"60e1c8719f09d6a04b69cff8","showtimeid":"60e2cd61ab5ef668319b72f2","movies":["60e272e5ac81947e79847c17","60e2bff5ac81947e79847c19"]}' 


post2:
	curl -XPOST http://localhost:4000/api/bookings/ -d '{"userid":"60e1c8e39f09d6a04b69cff9","showtimeid":"14501","movies":["The Ring","The Ring Two"]}' 	


build-docker:
	docker build . --tag renegmedal/cinema-nats-bookings:latest

push:
	docker push renegmedal/cinema-nats-bookings:latest
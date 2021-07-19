init-project:
	go mod init cinema-app/website
 

build:
	go build -race -o website main.go

run: build
	./website -serveraddr=localhost:8080 -s nats://localhost:4222 -users USERS -movies MOVIES -showtimes SHOWTIMES -bookings BOOKINGS


up:
	docker-compose up -d --build

run-docker:
	curl localhost:8080/

build-docker:
	docker build . --tag renegmedal/cinema-nats-website:latest

push:
	docker push renegmedal/cinema-nats-website:latest


#  
# $ make install-nats
# $ make install-nack
# $ make portforward-nats
# 
# $ make install-mongo-users
# $ make portforward-mongo-users
# $ cd users 
# $ make run
#
# $ cd website
# $ make run

# http://localhost:8080 
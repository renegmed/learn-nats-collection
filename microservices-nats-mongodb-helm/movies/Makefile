init-project:
	go mod init cinema-app/movies
 

build:
	go build -race -o movies cmd/app/* 

run: build
	./movies -serveraddr localhost:5000 -mongouri mongodb://localhost:27018 -enablecredentials=false -s nats://localhost:4222 -t=MOVIES -q qmovies


build-docker:
	docker build . --tag renegmedal/cinema-nats-movies:latest

push:
	docker push renegmedal/cinema-nats-movies:latest


# run:
# 	MONGODB_USERNAME=admin MONGODB_PASSWORD=pass123  go run cmd/app/*.go --mongoURI mongodb://localhost:27018

get-all:
	curl localhost:8282/api/movies/ -v

post1:
	curl -XPOST localhost:8282/api/movies/ -d '{"title":"The Thing", "director":"John Carpenter", "rating": 3.5, "createdon": "2012-05-25T00.00.00Z"}' -v 

post2:
	curl -XPOST localhost:8282/api/movies/ -d '{"title":"The Thing 2", "director":"John Carpenter", "rating": 4.5, "createdon": "2015-12-15T00.00.00Z"}' -v 

post3:
	curl -XPOST localhost:8282/api/movies/ -d '{"title":"The Mechanic", "director":"Peter Jackson", "rating": 4.0, "createdon": "2016-08-12T00.00.00Z"}' -v 
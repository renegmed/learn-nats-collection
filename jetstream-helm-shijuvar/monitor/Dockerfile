FROM golang:1.16.5-alpine

WORKDIR /src/monitor

COPY . .
 
RUN pwd
RUN ls -als  

RUN go mod vendor 

RUN go build -o monitor . 

CMD ["./monitor"]

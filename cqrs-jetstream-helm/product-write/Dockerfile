FROM golang:1.16.5-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
COPY main.go connector.go ./  

RUN go mod vendor
RUN go build -o product-publisher .
 

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
 
COPY --from=builder /app/product-publisher . 

RUN pwd
RUN ls -als /root

 
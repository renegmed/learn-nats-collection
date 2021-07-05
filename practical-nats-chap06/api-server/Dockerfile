FROM golang:1.16.5-alpine AS builder

WORKDIR /src 

COPY go.mod go.sum ./
COPY cmd cmd/
COPY kit kit/
COPY server server/ 


RUN go mod download
RUN go build -o api-server cmd/main.go
 

RUN pwd 
RUN ls -als 
# RUN ls -als /src



FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/ 
COPY --from=builder /src/api-server . 

RUN pwd
RUN ls -als /root 

# EXPOSE 8080
ENTRYPOINT ["/root/api-server"]


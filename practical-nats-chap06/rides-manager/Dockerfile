FROM golang:1.16.5-alpine AS builder

WORKDIR /src 

COPY go.mod go.sum ./
COPY cmd cmd/
COPY kit kit/
COPY server server/ 


RUN go mod download
RUN go build -o rides-manager cmd/main.go
 

RUN pwd 
RUN ls -als 
# RUN ls -als /src


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/ 
COPY --from=builder /src/rides-manager . 
 

# EXPOSE 8080
ENTRYPOINT ["/root/rides-manager"]

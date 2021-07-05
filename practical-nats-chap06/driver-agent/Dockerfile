 
FROM golang:1.16.5-alpine AS builder

WORKDIR /src 

COPY go.mod go.sum ./
COPY cmd cmd/
COPY kit kit/
COPY agent agent/ 


RUN go mod download
RUN go build -o driver-agent cmd/main.go
 

RUN pwd 
RUN ls -als 
# RUN ls -als /src


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/ 
COPY --from=builder /src/driver-agent . 

# EXPOSE 8080
ENTRYPOINT ["/root/driver-agent"]
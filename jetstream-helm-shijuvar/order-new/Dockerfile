FROM golang:1.16.5-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
COPY main.go ./  
COPY model/ model/

RUN go mod download
RUN go build -o order-new .
 

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/order-new . 

RUN pwd
RUN ls -als /root

# CMD ["./order-new","-s", "nats://localhost:4222"]
 
 
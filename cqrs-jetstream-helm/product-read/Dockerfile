FROM golang:1.16.5-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
COPY main.go connector.go ./  

RUN go mod vendor
RUN go build -o product-reader .
 

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/product-reader . 

RUN pwd
RUN ls -als /root

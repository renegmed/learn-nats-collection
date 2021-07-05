# golang alpine 1.15.6-alpine as base image
FROM golang:1.16.5-alpine AS builder
 
# create workspace
WORKDIR /app
COPY go.mod go.sum ./
 
COPY cmd cmd/
COPY pkg pkg/

# fetch dependancies
RUN go mod download
 

RUN pwd 
RUN ls -als

# build binary
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/cinema-bookings ./cmd/app
RUN go build -o movies cmd/app/*.go


# build a small image
FROM alpine:latest
LABEL language="golang"
 

WORKDIR /
COPY --from=builder /app/movies .

ENTRYPOINT ["./movies"]

 
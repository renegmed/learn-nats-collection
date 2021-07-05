# golang alpine 1.15.6-alpine as base image
FROM golang:1.16.5-alpine AS builder
# create appuser.
#RUN adduser -D -g '' elf
# create workspace
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY cmd/ ./cmd
COPY pkg/ ./pkg

# fetch dependancies
RUN go mod download
# RUN go mod verify

RUN pwd 
RUN ls -als

# build binary
#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/cinema-bookings ./cmd/app
RUN go build -o bookings cmd/app/*.go
#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o bookings ./cmd/app

# build a small image
FROM alpine:latest
LABEL language="golang"
 
WORKDIR /
COPY --from=builder /app/bookings .

ENTRYPOINT ["./bookings"]
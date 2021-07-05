FROM golang:1.16.5-alpine AS builder
# create appuser.
# RUN adduser -D -g '' elf
# create workspace
WORKDIR /app
COPY go.mod go.sum ./

COPY cmd cmd/
COPY pkg pkg/

# fetch dependancies
RUN go mod download
#RUN go mod verify

RUN pwd 
RUN ls -als

# build binary
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/cinema-bookings ./cmd/app
RUN go build -o users cmd/app/*.go 


# build a small image
FROM alpine:latest
LABEL language="golang"
#LABEL org.opencontainers.image.source cinema-app/bookings
# import the user and group files from the builder
#COPY --from=builder /etc/passwd /etc/passwd
# copy the static executable
#COPY --from=builder --chown=elf:1000 /go/bin/cinema-bookings /cinema-bookings
# use a non-root user
#USER elf
# run app

WORKDIR /
COPY --from=builder /app/users .

ENTRYPOINT ["./users"]



# # golang alpine 1.15.6-alpine as base image
# FROM golang:1.15.7-alpine AS builder
# # create appuser
# RUN adduser -D -g '' elf
# # create workspace
# WORKDIR /opt/app/
# COPY go.mod .
# COPY go.sum .
# # fetch dependancies
# RUN go mod download
# RUN go mod verify
# # copy the source code as the last step
# COPY . .
# # build binary
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/cinema-users ./cmd/app


# # build a small image
# FROM alpine:3.12.3
# LABEL language="golang"
# LABEL org.opencontainers.image.source cinema-app
# # import the user and group files from the builder
# COPY --from=builder /etc/passwd /etc/passwd
# # copy the static executable
# COPY --from=builder --chown=elf:1000 /go/bin/cinema-users /cinema-users
# # use a non-root user
# USER elf
# # run app
# ENTRYPOINT ["./cinema-users"]

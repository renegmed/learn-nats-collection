FROM golang:1.16.5-alpine AS builder
# create appuser.
# RUN adduser -D -g '' elf
# create workspace
WORKDIR /app
COPY go.mod go.sum ./

COPY main.go .
COPY web web/
COPY pkg pkg/
COPY ui ui/


RUN pwd 
RUN ls -als


# fetch dependancies
RUN go mod tidy
#RUN go mod verify

# build binary
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o website ./cmd/web

# RUN CGO_ENABLED=0 GOOS=linux go build -o website main.go
RUN go build -o website main.go 



# build a small image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
 

WORKDIR /
COPY --from=builder /app/website .
# COPY --from=builder /app/ui ui
 
# RUN pwd 
# RUN ls -als 
# # RUN ls -als ./ui
# # RUN ls -als ./ui/html
# # RUN ls -als ./ui/html/bookings

ENTRYPOINT ["./website"] 

FROM golang:1.16.5-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o cmds .


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/cmds .

ENTRYPOINT ["./cmds"]
                                
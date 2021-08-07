FROM golang:1.16.5-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o events .


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/events .

ENTRYPOINT ["./events"]
                                
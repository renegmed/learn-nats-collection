FROM golang:1.16.5-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o query .


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/query .

ENTRYPOINT ["./query"]
                                
FROM golang:1.16.5-alpine as builder

WORKDIR /app

COPY go.mod . 
COPY go.sum . 
COPY *.go ./  

RUN go mod download  

RUN go build -o subscriber . 

RUN pwd
RUN ls -als   

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /
 
COPY --from=builder /app/subscriber .

# ENTRYPOINT ["./subscriber"]
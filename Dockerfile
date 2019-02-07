FROM golang:1.11.5-alpine AS builder
WORKDIR /go/src/chadbot
ADD . /go/src/chadbot
RUN apk add --no-cache git \
    && go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chad .

FROM scratch  
WORKDIR /root/
COPY --from=builder /go/src/chadbot/chad .

ENTRYPOINT ["./chad"]  
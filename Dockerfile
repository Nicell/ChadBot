FROM golang:1.11.5-alpine
WORKDIR /go/src/chadbot
ADD . /go/src/chadbot
RUN apk add --no-cache ca-certificates ffmpeg gcc git libc-dev \
    && go get -d ./... \
    && go get -u github.com/bwmarrin/dca/cmd/dca
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chad .

ENTRYPOINT ["./chad"]

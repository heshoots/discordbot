FROM golang:1.9.3 as builder
RUN apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
    wget -O dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
    echo '31144e465e52ffbc0035248a10ddea61a09bf28b00784fd3fdd9882c8cbb2315 dep' | sha256sum -c - && \
    chmod +x dep && mv dep /usr/bin
WORKDIR /go/src/github.com/heshoots/discordbot
ADD Gopkg.lock .
ADD Gopkg.toml .
ADD main.go .
ADD models/ ./models
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.compiled=`date -u +.%Y%m%d.%H%M%S` -w" -o discordbot ./main.go

FROM alpine:3.7
RUN apk add --update ca-certificates
COPY --from=builder /go/src/github.com/heshoots/discordbot/discordbot /root/discordbot
WORKDIR /root
CMD ["./discordbot"]

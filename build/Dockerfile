FROM golang:1.11rc2 as builder
RUN apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y
WORKDIR /go/src/github.com/heshoots/discordbot
ENV GO111MODULE=on
ADD cmd cmd
ADD pkg pkg
ADD go.mod go.mod
ADD go.sum go.sum
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.compiled=`date -u +.%Y%m%d.%H%M%S` -w" -o discordbot cmd/discordbot/main.go

FROM alpine:3.7
RUN apk add --update ca-certificates
COPY --from=builder /go/src/github.com/heshoots/discordbot/discordbot /root/discordbot
WORKDIR /root
CMD ["./discordbot"]

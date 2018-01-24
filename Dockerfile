FROM golang:1.9.3
ADD main.go .
RUN go get -v github.com/heshoots/discordbot
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o discordbot ./main.go
CMD ["./discordbot"]

FROM golang

LABEL maintainer="Vlad Kampov <vladyslav.kampov@gmail.com>"
ENV GO111MODULE=on

WORKDIR $GOPATH/src/github.com/vladkampov/url-shortener-telegram-bot
ADD . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build main.go

EXPOSE 50051

CMD ["./main"]

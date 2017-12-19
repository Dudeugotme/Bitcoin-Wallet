
FROM golang:latest

RUN mkdir -p /go/src/app

COPY . /go/src/app

WORKDIR /go/src/app/cmd/bitcoin

RUN go-wrapper download

RUN go-wrapper install

CMD ["go-wrapper", "run"]

EXPOSE 8234
#FROM golang

#ADD . /go/src/Bitcoin-Wallet
#WORKDIR /go/src/Bitcoin-Wallet

#RUN go get ./...
#RUN go install ./...

#WORKDIR /go/bin/
FROM alpine

RUN mkdir -p ~/.spvwallet

COPY ./dist/bitcoin /usr/local/bin/bitcoin

RUN chmod +x /usr/local/bin/bitcoin

ENTRYPOINT []

CMD bitcoin

EXPOSE 8234
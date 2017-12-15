FROM golang:latest
RUN mkdir -p ~/.spvwallet

ADD ./dist/bitcoin /go/bin/bitcoin

RUN chmod +x /go/bin/bitcoin 

CMD /go/bin/bitcoin

EXPOSE 8234


#FROM alpine:latest

#RUN mkdir -p ~/.spvwallet

#RUN apk --no-cache add bash go git ca-certificates
#RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
#ENV GOPATH /go
#ENV PATH /go/bin:$PATH
#ADD ./dist/bitcoin /go/bin/bitcoin

#RUN chmod +x /go/bin/bitcoin 

#CMD /go/bin/bitcoin

#EXPOSE 8234
FROM golang:onbuild
RUN mkdir -p ~/.spvwallet

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app
RUN go build -o bitcoin main.go

CMD ["./bitcoin"]

EXPOSE 8234
FROM golang:onbuild
RUN mkdir -p ~/.spvwallet

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 
CMD ["/app/main"]

EXPOSE 8234
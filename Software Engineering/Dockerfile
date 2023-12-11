FROM golang:1.19.3

WORKDIR /go/src/mighty-saver-rabbit
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["mighty-saver-rabbit"]

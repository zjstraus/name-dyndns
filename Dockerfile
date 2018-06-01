FROM golang:alpine

WORKDIR /go/src/github.com/zjstraus/name-dyndns
COPY . .
RUN go install -v ./...

CMD ["name-dyndns"]

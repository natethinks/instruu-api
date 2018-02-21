FROM golang:1.8
WORKDIR /go/src/github.com/natethinks/instruu-api
ADD . .
RUN go get -v ./...
RUN go install -v ./...
WORKDIR /go/src/github.com/natethinks/instruu-api/cmd/instruu-api
CMD ["instruu-api"]

FROM golang:1.8
WORKDIR /go/src/github.com/natethinks/instruu-api
ADD . .
RUN pwd
RUN ls
RUN go get -v ./...
RUN go install -v ./...
CMD ["instruu-api"]

FROM golang:latest
ADD /go/src/github.com/natethinks/instruu-api .
WORKDIR /go/src/github.com/natethinks/instruu-api/cmd/instruu-api
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
CMD ["instruu-api"]

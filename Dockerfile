FROM golang:1.15

WORKDIR /go/src/github.com/theconcierge/recipes
COPY . .

RUN go get -d -v ./...
RUN go build 

CMD ["./recipes"]

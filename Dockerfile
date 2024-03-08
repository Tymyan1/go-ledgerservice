FROM golang:1.22
WORKDIR /app
COPY go.mod ./
COPY *.go ./

RUN go build -v -o /go-server .

EXPOSE 8080

CMD [ "/go-server" ]

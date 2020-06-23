FROM golang
WORKDIR /app

COPY . .

RUN go install -race /app
ENTRYPOINT /go/bin/message-api
EXPOSE 9090

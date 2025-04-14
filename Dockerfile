FROM golang:latest
RUN go version

WORKDIR /app

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-db.sh

RUN go mod download
RUN go build -o pvz ./cmd/pvz/main.go


CMD ["./pvz"]
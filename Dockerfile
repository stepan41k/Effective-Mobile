FROM golang:1.24.2

RUN go version

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x ./wait-for-postgres.sh

RUN go mod download
RUN go build -o profiles-library-service ./cmd/profiles/main.go
CMD ["./profiles-library-service"]
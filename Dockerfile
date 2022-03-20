FROM golang:1.17

WORKDIR /src/

COPY . /src/

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

RUN go mod download;

RUN go build -o scanner ./cmd/scanner/main.go

CMD [ "./scanner" ]


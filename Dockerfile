FROM golang:1.17-buster as build

ENV GOPATH=/

WORKDIR /src/

COPY . /src/

RUN go mod download

RUN go build -o scanner ./cmd/scanner/main.go

CMD [ "./scanner" ]


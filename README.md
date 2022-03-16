# Scanner

## Description

TG_scanner is an application which can you use to parse messages from your groups in telegram
Application will always create a folder with title data where you can get result and folder logs with file all.log for logs

## Technology

Go,GoTD,Logrus,GoDotenv

## Usage

Starting an application locally:

```bash
go mod download
make start #Or you can use go run ./cmd/scanner/main.go
```

Starting in docker:

```bash
docker build -t scanner .
docker run --name tg_scanner scanner
```

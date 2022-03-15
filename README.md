# Scanner

## Description

TG_scanner is an application which can you use to parse messages from your groups in telegram
Application will always create a folder with title data where you can get result and folder logs with file all.log for logs

## Technology

Go,GoTD,Logrus,GoDotenv

## Usage

Before start using this application you should create .env files with this fields

APP_ID=Telegram app id

APP_HASH=Telegram app hash

PHONE=Your mobile phone which you are going to use

PASSWORD=Password to telegram

LIMIT=Limit for getting messages from telegram history.1000 is the most perfect value


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

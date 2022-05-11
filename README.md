# Telegram message scanner

## Description

tg_scanner is an application which can parse question and replies from you telegram groups

Application will always create a few dirs:
  - logs - for logs file
  - images - for saving and getting user images before save to firebase
  - data - for saving and getting messages and replies from telegram before save to PostgreSQL

## Technology

Go, GoTD, PostgreSQL, Testify, go-sqlmock, firebase-admin-go

## Before start

Please create a dir "configs" with file ".config.env" which have this fields:
- APP_ID = Telegram app id
- APP_HASH = Telegram app hash
- PHONE = Telegram phone number 
- PASSWORD = Password to telegram [You must have it in telegram]
- POSTGRES_USER = Postgresql username
- POSTGRES_PASSWORD = Postgresql password
- POSTGRES_DB = Postgresql database name
- PROJECT_ID = Project id from firebase
- STORAGE_BUCKET = Storage bucket name from firebase
- SECRET_PATH = Path to your secret key from firebase

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

Starting with docker-compose:

```bash
 docker-compose build

 docker-compose up #After it enter code which telegram send to you
```

Running tests


```bash
 #Before start run: make mock

 make test
```

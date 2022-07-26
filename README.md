# tg_scanner

## Description

tg_scanner is an application which can parse question and replies from your telegram groups

Application will always create a few dirs:

- logs - for logs file
- images - for saving and getting user, message images before save to firebase
- data - for saving and getting data from telegram before save to PostgreSQL

## Technology

Go, GoTD, PostgreSQL, Testify, go-sqlmock, firebase-admin-go, Redis


## Before start

Please create a dir "configs" with file ".config.env" which have this fields:

### Telegram:

- APP_ID = Telegram app id
- APP_HASH = Telegram app hash
- PHONE = Telegram phone number
- PASSWORD = Password to telegram

### PostgreSQL:

- DATABASE_URL = You can also use it for describe PostgreSQL path
- MIGRATIONS_PATH = Path to file with migrations

### Firebase:

- PROJECT_ID = Project id from firebase
- STORAGE_BUCKET = Storage bucket name from firebase
- SECRET_PATH = Path to your secret key from firebase

### Logger:

- LOG_LEVEL = Log level which logger should handler

### Redis:

- REDIS_ADDR = Redis address
- REDIS_PASSWORD = Password for redis

## Usage

Starting an application locally:

```bash
 $ go mod download

 $ make start #Or you can use go run ./cmd/scanner/main.go
```

Starting in docker:

```bash
 $ make docker

 $ docker run --name telegram_scanner tg_scanner

```

Running tests

```bash
 $ make mock # Run it if mocks folder is not exist

 $ make test
```

Watch demo version:

[Telegram Overflow](https://telegram-overflow.herokuapp.com/)


# tg_scanner

tg_scanner is an application which can parse question and replies from your telegram groups

Application will always create a few dirs:
- logs - for logs file
- images - for saving and getting user, message images before save to firebase
- data - for saving and getting data from telegram before sent to message queue

## Tech Stack

**Main:** 
- gotd/td
- apache kafka 

**Store:**
- redis
- firebase storage

## Features

- Get channel, user, message reply info from telegram
- Save channel images, user images, images from replies and messages into Firebase Storage
- Microservice communication with Apache Kafka
- Caching messages with Redis


## Environment Variables

To run this project, you will need to add the following environment variables to your ".config.env" file which locate in "config" folder:

#### Telegram
- `APP_ID`- Telegram app id
- `APP_HASH`- Telegram app hash
- `PHONE` - Your phone number from telegram account
- `PASSWORD` - Your password to telegram

#### Firebase
- `PROJECT_ID` - Firebase project id
- `STORAGE_BUCKET` - Firebase storage bucket
- `SECRET_PATH` - Path to firebase secret file

#### Logger
- `LOG_LEVEL` - Level which logger will process

#### Redis
- `REDIS_ADDR` - Redis address
- `REDIS_PASSWORD` - Redis password

#### Kafka
- `KAFKA_ADDR` - Kafka address
## Run Locally

Clone the project

```bash
  git clone git@github.com:VladPetriv/tg_scanner.git
```

Go to the project directory

```bash
  cd tg_scanner
```

Install dependencies

```bash
  go mod download
```

Start the application:

Before start you must create 2 topic in kafka: "channels.get", "messages.get"

```bash
  # Make sure that Apache Kafka is running
  make run # Or you can use "go run ./cmd/scanner/main.go"
```

Start the application with docker:

```bash
  make docker # Build docker image

  docker run --name telegram_scanner tg_scanner
```

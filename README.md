
# tg_scanner

tg_scanner is an application which can parse question and replies from your telegram groups.

Application will always create a few dirs:
- images - for saving and getting user, group, message and reply images before save to firebase. (The images will auto remove after saving to the Firebase)
- data - for saving and getting data from telegram before sent to message queue.

## Tech Stack

**Main:** 
- gotd/td
- apache kafka 

**Store:**
- redis
- firebase storage

## Features

- Get groups, user, message reply info from telegram.
- Save groups, user, replies and messages images into Firebase.
- Microservice communication with Apache Kafka.
- Caching messages and groups with Redis.
- Auto deleting log files after 22 days


## Environment Variables

To run this project, you will need to add the following environment variables to your `config/.config.env` file:

#### Telegram
- `APP_ID`- Telegram app id
- `APP_HASH`- Telegram app hash
- `PHONE` - Your phone number from telegram account
- `PASSWORD` - Your cloud password to telegram

#### Firebase
- `PROJECT_ID` - Firebase project id
- `STORAGE_BUCKET` - Firebase storage bucket
- `SECRET_PATH` - Path to firebase secret file

#### Logger
- `LOG_LEVEL` - Level which logger will process
- `LOG_FILENAME` - Filepath where logger will save logs
#### Redis
- `REDIS_ADDR` - Redis address
- `REDIS_PASSWORD` - Redis password

#### Kafka
- `KAFKA_ADDR` - Kafka address
## Run

Clone the project

```bash
  git clone git@github.com:VladPetriv/tg_scanner.git
```

Go to the project directory

```bash
  cd tg_scanner
```

Install dependencies

```go
  go mod download
```

Start the application:

Before start you must create 2 topic in kafka: `groups` and `messages`

```bash
  make run 
```

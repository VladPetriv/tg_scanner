package main

import (
	"context"
	"log"
	"time"

	"example.com/test/m/config"
	"example.com/test/m/internal/auth"
	"example.com/test/m/internal/channel"
	"example.com/test/m/internal/file"
	"example.com/test/m/internal/filter"
	"example.com/test/m/internal/message"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)

	}

	client, err := telegram.ClientFromEnvironment(telegram.Options{})
	if err != nil {
		log.Fatalf("ERROR_WHILE_CREATING_CLIENT:%s", err)
	}

	//Create API
	api := client.API()

	//Create context
	ctx := context.Background()

	if err := client.Run(ctx, func(ctx context.Context) error {

		for {
			log.Println("Start")

			//Authorization to telegram
			_, err := auth.Login(ctx, *client, *cfg)
			if err != nil {
				return err
			}

			accessHash, err := channel.GetAccessHash(ctx, "nodejs_ru", api)
			if err != nil {
				return err
			}

			data, err := channel.GetChannelHistory(ctx, api, tg.InputPeerChannel{
				ChannelID:  1041204341,
				AccessHash: accessHash,
			}, 5)
			if err != nil {
				return err
			}

			modifiedData, _ := data.AsModified()

			messages := message.GetMessagesFromTelegram(ctx, modifiedData, api, &tg.InputPeerChannel{
				ChannelID:  1041204341,
				AccessHash: accessHash,
			})

			messagesFromFile, err := file.GetMessagesFromFile("message.json")
			if err != nil {
				return err
			}

			for _, message := range messages {
				msg, ok := filter.FilterMessages(&message)
				if !ok {
					continue
				}
				messagesFromFile = append(messagesFromFile, *msg)
			}

			result := filter.RemoveDuplicateByMessage(messagesFromFile)

			err = file.WriteMessagesToFile(result, "message.json")
			if err != nil {
				return err
			}

			log.Println("Completed without errors")
			time.Sleep(time.Second)
			log.Println("Wait 30s and do request again")
			time.Sleep(time.Second * 30)
		}
	}); err != nil {
		panic(err)
	}

}

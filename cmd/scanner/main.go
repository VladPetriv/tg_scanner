package main

import (
	"context"
	"fmt"
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

		//Authorization to telegram
		_, err = auth.Login(ctx, *client, *cfg)
		if err != nil {
			log.Fatal(err)
		}
		groups, err := channel.GetAllGroups(ctx, api)
		if err != nil {
			return err
		}

		file.CreateFiles(groups)
		for _, group := range groups {
			go GetResult(ctx, group, api)
		}
		fmt.Scanln()
		return nil
	}); err != nil {
		panic(err)
	}

}

func GetResult(ctx context.Context, group channel.Group, api *tg.Client) error {
	fileName := fmt.Sprintf("%s.json", group.Username)
	for {
		log.Printf("Start with %s", group.Username)

		data, err := channel.GetChannelHistory(ctx, api, tg.InputPeerChannel{
			ChannelID:  int64(group.ID),
			AccessHash: int64(group.AccessHash),
		}, 5)
		if err != nil {
			return err
		}

		modifiedData, _ := data.AsModified()

		messages := message.GetMessagesFromTelegram(ctx, modifiedData, api, &tg.InputPeerChannel{
			ChannelID:  int64(group.ID),
			AccessHash: int64(group.AccessHash),
		})

		messagesFromFile, err := file.GetMessagesFromFile(fileName)
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

		err = file.WriteMessagesToFile(result, fileName)
		if err != nil {
			return err
		}

		log.Printf("Completed without errors [%s]", group.Username)
		time.Sleep(time.Second)
		log.Printf("Wait 30s and do request again[%s]", group.Username)
		time.Sleep(time.Second * 30)
	}
}

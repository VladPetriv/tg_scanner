package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/auth"
	"github.com/VladPetriv/tg_scanner/internal/channel"
	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/filter"
	"github.com/VladPetriv/tg_scanner/internal/message"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func main() {
	//Initialize config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	//Createting dir for data
	err = file.CreateDir()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup

	//Create new client
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
		user, err := auth.Login(ctx, client, cfg)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(2)
		//Get user data
		u, _ := user.GetUser().AsNotEmpty()

		//Getting incoming messages
		go GetNewMessage(ctx, u, api, &wg)

		//Getting all groups
		groups, err := channel.GetAllGroups(ctx, api)
		if err != nil {
			return err
		}

		//Create files
		file.CreateFilesForGroups(groups)

		//Getting group history
		for _, group := range groups {
			go GetFromHistory(ctx, group, api, cfg, &wg)
		}
		wg.Wait()
		return nil
	}); err != nil {
		panic(err)
	}

}

func GetFromHistory(ctx context.Context, group channel.Group, api *tg.Client, cfg *config.Config, wg *sync.WaitGroup) error {
	defer wg.Done()
	fileName := fmt.Sprintf("./data/%s.json", group.Username)
	for {
		log.Printf("Start with %s", group.Username)

		data, err := channel.GetChannelHistory(ctx, cfg.Limit, tg.InputPeerChannel{
			ChannelID:  int64(group.ID),
			AccessHash: int64(group.AccessHash),
		}, api)
		if err != nil {
			return err
		}

		modifiedData, _ := data.AsModified()

		messages := message.GetMessagesFromTelegram(ctx, modifiedData, &tg.InputPeerChannel{
			ChannelID:  int64(group.ID),
			AccessHash: int64(group.AccessHash),
		}, api)

		messagesFromFile, err := file.GetMessagesFromFile(fileName)
		if err != nil {
			return err
		}

		for _, m := range messages {
			msg, ok := filter.FilterMessages(&m)
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
		time.Sleep(time.Second * 30)
	}
}

func GetNewMessage(ctx context.Context, user *tg.User, api *tg.Client, wg *sync.WaitGroup) error {
	defer wg.Done()
	fileName := "incoming.json"
	path := fmt.Sprintf("./data/%s", fileName)
	err := file.CreateFileForIncoming()
	if err != nil {
		return err
	}

	for {
		log.Println("Start getting new message)")

		messagesFromFile, err := file.GetMessagesFromFile(path)
		if err != nil {
			return err
		}

		incomingMessage, err := message.GetIncomingMessages(ctx, user, api)
		if err != nil {
			return err
		}

		for _, m := range incomingMessage {
			msg, ok := filter.FilterMessages(&m)
			if !ok {
				continue
			}
			messagesFromFile = append(messagesFromFile, *msg)
		}
		result := filter.RemoveDuplicateByMessage(messagesFromFile)

		err = file.WriteMessagesToFile(result, path)
		if err != nil {
			return err
		}

		log.Println("Completed without errors)")
		time.Sleep(time.Minute)
	}
}

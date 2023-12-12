package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Timestamp)
		fmt.Println(event.Event)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf(err.Error())
	}

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	witClient := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}

	// go printCommandEvents(bot.CommandEvents())

	bot.Command("query for my bot - <query>", &slacker.CommandDefinition{
		Description: "ai chat bot",
		Examples:    []string{"what is the capital of India"},
		Handler: func(ctx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("query")

			// call wit with the user input to construct query which wolfram can understand
			res, err := witClient.Parse(&witai.MessageRequest{
				Query: query,
			})
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// convert response to understandable json format to fetch required value for processing
			data, err := json.Marshal(res)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			rough := string(data[:])

			// use gjson to deconstruct the response (gjson.ResponseToUnderstandableJson(res, "$wolframquer...") which is very nested json
			val := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			stringVal := val.String()

			// call wolfram with the query result decoded from the wit response
			receivedAnswer := "NA"
			receivedAnswer, err = wolframClient.GetSpokentAnswerQuery(stringVal, wolfram.Metric, 1000)
			if err != nil {
				log.Fatal(err)
			}

			response.Reply(receivedAnswer, slacker.WithThreadReply(true))
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Fatal(bot.Listen(ctx))

	// File upload via go server to slack channel
	fileBot := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	channelArr := []string{os.Getenv("CHANNEL_ID")}
	files := []string{"sample.csv"}

	for i := range files {
		params := slack.FileUploadParameters{
			Channels: channelArr,
			File:     files[i],
		}

		file, err := fileBot.UploadFile(params)
		if err != nil {
			log.Print(err.Error())
		}
		log.Printf("URL: %v; Name: %v", file.URLPrivate, file.Name)
	}
}

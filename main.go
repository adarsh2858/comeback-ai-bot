package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
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
	fmt.Println("Hello world")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	witClient := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	// wolframClient, _ := wit.NewClient(os.Getenv("WOLFRAM_APP_TOKEN"))

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

			data, err := json.Marshal(res)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			rough := string(data[:])

			// use gjson to deconstruct the response (gjson.ResponseToUnderstandableJson(res, "$wolframquery.wolframquery").0.value) which is  very nested json
			val := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			fmt.Println(val)

			// call wolfram with the query result decoded from the wit response
			response.Reply("Hello world", slacker.WithThreadReply(true))
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Fatal(bot.Listen(ctx))
}

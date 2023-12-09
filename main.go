package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

// func printCommandEvents(analyticsChannel <- chan *slacker.CommandEvent) {
// for event := range analyticsChannel{
// 	fmt.Println(event.Command)
// 	fmt.Println(event.Parameters)
// 	fmt.Println(event.Timestamp)
// 	fmt.Println(event.Event)}
// }

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println("Hello world")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	// witClient, _ := wit.NewClient(os.Getenv("WIT_APP_TOKEN"))
	// wolframClient, _ := wit.NewClient(os.Getenv("WOLFRAM_APP_TOKEN"))

	// go printCommandEvents(bot.CommandEvents())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Command("query for my bot - <query>", &slacker.CommandDefinition{
		Description: "ai chat bot",
		Examples:    []string{"what is the capital of India"},
		Handler: func(ctx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("query")
			fmt.Println(query)

			// call wit with the user input to construct query which wolfram understand
			// res, err := witClient.CommunicateForProcessing(witClient.QueryStruct{
			// 	Query: query
			// })
			// if err != nil {
			// 	fmt.Println(err.Error())
			// 	return
			// }

			// use gjson to deconstruct the response (gjson.ResponseToUnderstandableJson(res, "$wolframquery.wolframquery").0.value) which is  very nested json

			// call wolfram with the query result decoded from the wit response
			response.Reply("Hello world", slacker.WithThreadReply(true))
		},
	})

	log.Fatal(bot.Listen(ctx))
}

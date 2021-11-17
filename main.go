package main

import (
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/jiramot/go-linebot/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	bot, err := linebot.New(config.LineConfig.ChannelSecret, config.LineConfig.ChannelAccessToken)
	if err != nil {
		log.Fatal(err)
	}

	e.POST("/", func(c echo.Context) error {
		events, err := bot.ParseRequest(c.Request())
		if err != nil {
			log.Panic(err)
			return c.String(http.StatusBadRequest, "")
		}

		for _, event := range events {
			fmt.Printf("UserID: %s\n", event.Source.UserID)
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					replyMsg := linebot.NewTextMessage(message.Text)
					if strings.EqualFold(message.Text, "location") {
						locationAction := linebot.NewLocationAction("Share location")
						quickReplyButton := linebot.NewQuickReplyButton("", locationAction)
						quickReplyItem := linebot.NewQuickReplyItems(quickReplyButton)
						replyMsg.WithQuickReplies(quickReplyItem)
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, replyMsg).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.LocationMessage:
					cols := []*linebot.CarouselColumn{}
					googleAction := linebot.NewURIAction("action", "https://google.com")
					for i := 0; i < 10; i++ {
						col := linebot.CarouselColumn{
							Title:         fmt.Sprintf("Title %d", i),
							Text:          fmt.Sprintf("text %d", i),
							Actions:       []linebot.TemplateAction{googleAction},
							DefaultAction: googleAction,
						}
						cols = append(cols, &col)
					}
					template := linebot.NewCarouselTemplate(cols...)
					replyMsg := linebot.NewTemplateMessage("branch", template)
					if _, err = bot.ReplyMessage(event.ReplyToken, replyMsg).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.ServerConfig.Port)))
}
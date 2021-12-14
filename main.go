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

					if strings.EqualFold(message.Text, "location") {
						urlAction := linebot.NewURIAction("Share location", "https://line.me/R/nv/location/")
						template := linebot.NewButtonsTemplate("", "Title", "text", urlAction)
						replyMsg := linebot.NewTemplateMessage("branch", template)
						if _, err = bot.ReplyMessage(event.ReplyToken, replyMsg).Do(); err != nil {
							log.Print(err)
						}
					} else {
						replyMsg := linebot.NewTextMessage(message.Text)
						locationAction := linebot.NewLocationAction("Share location")
						quickReplyButton := linebot.NewQuickReplyButton("", locationAction)
						quickReplyItem := linebot.NewQuickReplyItems(quickReplyButton)
						replyMsg.WithQuickReplies(quickReplyItem)
						if _, err = bot.ReplyMessage(event.ReplyToken, replyMsg).Do(); err != nil {
							log.Print(err)
						}
					}
				case *linebot.LocationMessage:
					// cols := []*linebot.CarouselColumn{}
					// googleAction := linebot.NewURIAction("action", "https://google.com")
					// for i := 0; i < 10; i++ {
					// 	col := linebot.CarouselColumn{
					// 		Title:         fmt.Sprintf("Title %d", i),
					// 		Text:          fmt.Sprintf("text %d", i),
					// 		Actions:       []linebot.TemplateAction{googleAction},
					// 		DefaultAction: googleAction,
					// 	}
					// 	cols = append(cols, &col)
					// }
					// template := linebot.NewCarouselTemplate(cols...)
					// replyMsg := linebot.NewTemplateMessage("branch", template)
					// if _, err = bot.ReplyMessage(event.ReplyToken, replyMsg).Do(); err != nil {
					// 	log.Print(err)
					// }

					var contents []*linebot.BubbleContainer
					for i := 0; i < 10; i++ {
						flex1 := 1
						flex5 := 5
						bubble := &linebot.BubbleContainer{
							Hero: &linebot.ImageComponent{
								URL:         "https://lh3.ggpht.com/p/AF1QipMb5o6KHEhTwcPbmquKX3BQg0cCFlM05PAaJ2qS=s1024",
								Size:        linebot.FlexImageSizeTypeFull,
								AspectRatio: linebot.FlexImageAspectRatioType20to13,
								AspectMode:  linebot.FlexImageAspectModeTypeCover,
							},
							Body: &linebot.BoxComponent{
								Layout: linebot.FlexBoxLayoutTypeVertical,
								Contents: []linebot.FlexComponent{
									&linebot.TextComponent{
										Text:   "CJ Express",
										Weight: linebot.FlexTextWeightTypeBold,
										Size:   linebot.FlexTextSizeTypeXl,
									},
									&linebot.BoxComponent{
										Layout:  linebot.FlexBoxLayoutTypeVertical,
										Spacing: linebot.FlexComponentSpacingTypeSm,
										Contents: []linebot.FlexComponent{
											&linebot.TextComponent{
												Text:  "Place",
												Color: "#aaaaaa",
												Size:  linebot.FlexTextSizeTypeSm,
												Flex:  &flex1,
											},
											&linebot.TextComponent{
												Text:   "393 Si Lom, Silom, Bang Rak, Bangkok 10500",
												Wrap:   true,
												Color:  "#666666",
												Weight: linebot.FlexTextWeightTypeBold,
												Size:   linebot.FlexTextSizeTypeXl,
												Flex:   &flex5,
											},
										},
									},
									&linebot.BoxComponent{
										Layout:  linebot.FlexBoxLayoutTypeBaseline,
										Spacing: linebot.FlexComponentSpacingTypeSm,
										Contents: []linebot.FlexComponent{
											&linebot.TextComponent{
												Text:  "Time",
												Color: "#aaaaaa",
												Size:  linebot.FlexTextSizeTypeSm,
												Flex:  &flex1,
											},
											&linebot.TextComponent{
												Text:  "9:00 - 18:00",
												Color: "#666666",
												Size:  linebot.FlexTextSizeTypeSm,
												Flex:  &flex5,
											},
										},
									},
								},
							},
						}
						contents = append(contents, bubble)
					}
					carouselContainer := linebot.CarouselContainer{
						Type:     linebot.FlexContainerTypeCarousel,
						Contents: contents,
					}
					flex := linebot.NewFlexMessage("alt", &carouselContainer)
					if _, err = bot.ReplyMessage(event.ReplyToken, flex).Do(); err != nil {
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

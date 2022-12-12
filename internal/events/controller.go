package events

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Controller struct {
	Bot   *linebot.Client
	store *LocalStore
}

func NewController(bot *linebot.Client) *Controller {
	return &Controller{Bot: bot,
		store: NewLocalStore(),
	}
}

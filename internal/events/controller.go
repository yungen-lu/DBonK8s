package events

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/yungen-lu/TOC-Project-2022/internal/client"
)

type Controller struct {
	Bot       *linebot.Client
	store     *LocalStore
	k8sclient *client.K8sClient
}

func NewController(bot *linebot.Client, cl *client.K8sClient) *Controller {
	return &Controller{Bot: bot,
		store:     NewLocalStore(),
		k8sclient: cl,
	}
}

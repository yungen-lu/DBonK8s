package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/yungen-lu/TOC-Project-2022/internal/client"
)

func NewRouter(mux *chi.Mux, bot *linebot.Client, cl *client.K8sClient) {
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Recoverer)
	mux.Group(newWebHookHandler(bot, cl))

}

package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func NewRouter(mux *chi.Mux, bot *linebot.Client) {
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Recoverer)
	mux.Group(newWebHookHandler(bot))
	// mux.Post("/webhook", )

}

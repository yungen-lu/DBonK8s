package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/yungen-lu/TOC-Project-2022/internal/client"
	"github.com/yungen-lu/TOC-Project-2022/internal/events"
)

type webhookHandler struct {
	eventController *events.Controller
	k8sclient       *client.K8sClient
}

func newWebHookHandler(bot *linebot.Client, cl *client.K8sClient) func(r chi.Router) {
	whh := &webhookHandler{
		eventController: events.NewController(bot, cl),
		k8sclient:       cl,
	}
	// bot.ParseRequest()
	return func(r chi.Router) {
		r.Use(linebotMiddleWare(bot))
		r.Post("/webhook", whh.handle)

	}
}
func linebotMiddleWare(bot *linebot.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			event, err := bot.ParseRequest(r)
			if err != nil {
				log.Warn(err.Error())
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), "events", event))
			next.ServeHTTP(w, r)
		})

	}
}

// EventTypeMessage           EventType = "message"
// EventTypeFollow            EventType = "follow"
// EventTypeUnfollow          EventType = "unfollow"
// EventTypeJoin              EventType = "join"
// EventTypeLeave             EventType = "leave"
// EventTypeMemberJoined      EventType = "memberJoined"
// EventTypeMemberLeft        EventType = "memberLeft"
// EventTypePostback          EventType = "postback"
// EventTypeBeacon            EventType = "beacon"
// EventTypeAccountLink       EventType = "accountLink"
// EventTypeThings            EventType = "things"
// EventTypeUnsend            EventType = "unsend"
// EventTypeVideoPlayComplete EventType = "videoPlayComplete"
func (whh *webhookHandler) handle(w http.ResponseWriter, r *http.Request) {
	linebotevents, ok := r.Context().Value("events").([]*linebot.Event)
	if !ok {
		log.Warn("can't find event in context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, event := range linebotevents {
		log.WithField("event", event.Type).Info("request")
		var err error
		switch event.Type {
		case linebot.EventTypeMessage:
			err = whh.eventController.HandleEventTypeMessage(event)
		case linebot.EventTypePostback:
			err = whh.eventController.HandleEventTypePostBack(event)
		case linebot.EventTypeFollow:
			err = whh.eventController.HandleEventTypeFollow(event)
		case linebot.EventTypeUnfollow:
			// do nothing
		default:
			err = fmt.Errorf("no event hanlder for: %s\n", event.Type)
		}
		if err != nil {
			log.Warn(err.Error())
			_, err := whh.eventController.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).WithContext(r.Context()).Do()
			if err != nil {
				log.Warn(err.Error())
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

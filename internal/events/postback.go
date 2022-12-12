package events

import (
	"net/url"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (c *Controller) HandleEventTypePostBack(event *linebot.Event) error {
	data := event.Postback.Data
	replyToken := event.ReplyToken
	q, err := url.ParseQuery(data)
	if err != nil {
		return err
	}
	u, err := c.store.GetOrCreateUser(event.Source.UserID, c)
	if err != nil {
		return err
	}
	action := q.Get("action")
	dbname := q.Get("dbname")
	namespace := q.Get("ns")
	switch action {
	case "info":
		return u.FSM.Fire(InfoEvent, replyToken, dbname, namespace)
	case "delete":
		return u.FSM.Fire(StopEvent, replyToken, dbname, namespace)
	}
	return nil
}

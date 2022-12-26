package events

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
)

func (c *Controller) HandleEventTypeFollow(event *linebot.Event) error {
	replyToken := event.ReplyToken
	u, err := c.store.GetOrCreateUser(event.Source.UserID, c)
	if err != nil {
		return err
	}
	followMessage := `Thanks for using DBonK8s bot
Available commands are as follows
1. config [-upt]
2. list [-an]
3. info [-dn]
4. stop [-dn]
5. create [-dn]
6. myinfo [-a]
7. fsm
add -h flag on each command to get more info`
	_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(followMessage)).Do()
	if err != nil {
		log.Warn(err.Error())
	}
	return err
}

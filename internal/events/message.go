package events

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
)

// MessageTypeText     MessageType = "text"
// MessageTypeImage    MessageType = "image"
// MessageTypeVideo    MessageType = "video"
// MessageTypeAudio    MessageType = "audio"
// MessageTypeFile     MessageType = "file"
// MessageTypeLocation MessageType = "location"
// MessageTypeSticker  MessageType = "sticker"
// MessageTypeTemplate MessageType = "template"
// MessageTypeImagemap MessageType = "imagemap"
// MessageTypeFlex     MessageType = "flex"
type configOpts struct {
	UserName   string `short:"u" long:"username" description:"configure username"`
	PassWord   string `short:"p" long:"password" description:"configure password"`
	AdminToken string `short:"t" long:"token" description:"admin token"`
}

type listOpts struct {
	All       bool   `short:"a" long:"all" description:"list all namespace instance"`
	Namespace string `short:"n" long:"namespace" description:"database namespace"`
}

type infoOpts struct {
	DBName    string `short:"d" long:"dbname" description:"database name" required:"true"`
	Namespace string `short:"n" long:"namespace" description:"database namespace"`
}

type stopOpts infoOpts

type createOpts infoOpts

func (c *Controller) HandleEventTypeMessage(event *linebot.Event) error {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		log.Info("handling text message")
		return c.handleText(message, event.ReplyToken, event.Source)
	default:
		return fmt.Errorf("can't handle message type: %s\n", event.Message.Type())
	}
}
func (c *Controller) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	if source.UserID == "" {
		return errors.New("can't get user id")
	}
	// u, err := c.store.GetUser(source.UserID)
	u, err := c.store.GetOrCreateUser(source.UserID, c)
	if err != nil {
		return err
	}
	args := strings.Fields(strings.ToLower(message.Text))
	if len(args) < 1 {
		return errors.New("args len is lower then 1")
	}

	cmd := args[0]
	args = args[1:]
	print(cmd, args)

	// var err error
	switch cmd {
	case "config":
		var opts configOpts
		args, err = flags.ParseArgs(&opts, args)
		if err != nil {
			return fmt.Errorf("config opts: %s", err.Error())
		}
		return u.FSM.Fire(ConfigEvent, replyToken, opts.UserName, opts.PassWord, opts.AdminToken)
	case "list":
		var opts listOpts
		args, err = flags.ParseArgs(&opts, args)
		if err != nil {
			return fmt.Errorf("list opts: %s", err.Error())
		}
		return u.FSM.Fire(ListEvent, replyToken, opts.All, opts.Namespace)
	case "info":
		var opts infoOpts
		args, err = flags.ParseArgs(&opts, args)
		if err != nil {
			return fmt.Errorf("info opts: %s", err.Error())
		}
		return u.FSM.Fire(InfoEvent, replyToken, opts.DBName, opts.Namespace)
	case "stop":
		var opts stopOpts
		args, err = flags.ParseArgs(&opts, args)
		if err != nil {
			return fmt.Errorf("stop opts: %s", err.Error())
		}
		return u.FSM.Fire(StopEvent, replyToken, opts.DBName, opts.Namespace)
	case "create":
		var opts createOpts
		args, err = flags.ParseArgs(&opts, args)
		if err != nil {
			return fmt.Errorf("create opts: %s", err.Error())
		}
		return u.FSM.Fire(CreateEvent, replyToken, opts.DBName, opts.Namespace)
	default:
		return errors.New("command not found")
	}

	// _, err := c.bot.ReplyMessage(replyToken, linebot.NewTextMessage("echo: "+message.Text)).Do()
}

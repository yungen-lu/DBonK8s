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

type createOpts struct {
	DBName    string `short:"d" long:"dbname" description:"database name" required:"true"`
	DBType    string `short:"t" long:"type" description:"database type" required:"true" choice:"postgres" choice:"mysql" choice:"redis" choice:"mongodb"`
	Namespace string `short:"n" long:"namespace" description:"database namespace"`
}

type userInfoOpts struct {
	All bool `short:"a" long:"all" description:"list all user"`
}

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
	log.WithFields(log.Fields{"user": u.UserID, "cmd": cmd, "args": args}).Info("user sended command")

	switch cmd {
	case "config":
		var opts configOpts
		err = Parse("config", &opts, args)
		if err != nil {
			return err
		}
		return u.FSM.Fire(ConfigEvent, replyToken, opts.UserName, opts.PassWord, opts.AdminToken)
	case "list":
		var opts listOpts
		err = Parse("list", &opts, args)
		if err != nil {
			return err
		}
		return u.FSM.Fire(ListEvent, replyToken, opts.All, opts.Namespace)
	case "info":
		var opts infoOpts
		err = Parse("info", &opts, args)
		if err != nil {
			return err
		}
		return u.FSM.Fire(InfoEvent, replyToken, opts.DBName, opts.Namespace)
	case "stop":
		var opts stopOpts
		err = Parse("stop", &opts, args)
		if err != nil {
			return err
		}
		return u.FSM.Fire(StopEvent, replyToken, opts.DBName, opts.Namespace)
	case "create":
		var opts createOpts
		err = Parse("create", &opts, args)
		if err != nil {
			return err
		}
		return u.FSM.Fire(CreateEvent, replyToken, opts.DBName, opts.DBType, opts.Namespace)
	case "myinfo":
		var opts userInfoOpts
		err = Parse("myinfo", &opts, args)
		if err != nil {
			return err
		}
		return u.FSM.Fire(UserInfoEvent, replyToken, opts.All)
	case "back":
		return u.FSM.Fire(BackEvent)
	case "fsm":
		return u.FSM.Fire(FSMEvent, replyToken)
	default:
		return errors.New("command not found")
	}

}
func Parse(name string, data interface{}, args []string) error {
	parser := flags.NewNamedParser(name, flags.Default)
	_, err := parser.AddGroup("Application Options", "", data)
	if err != nil {
		return err
	}
	_, err = parser.ParseArgs(args)
	if err != nil {
		return err
	}
	return nil
}

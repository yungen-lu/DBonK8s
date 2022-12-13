package events

import (
	"context"
	"reflect"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/qmuntal/stateless"
	log "github.com/sirupsen/logrus"
)

type User struct {
	UserID   string
	UserName string
	PassWord string
	IsAdmin  bool
	FSM      *stateless.StateMachine
	Con      *Controller
}

const (
	UserState        = "user"
	ListStateUser    = "list[user]"
	ConfigState      = "config"
	InfoStateUser    = "info[user]"
	StopStateUser    = "stop[user]"
	CreateStateUser  = "create[user]"
	AdminState       = "admin"
	ListStateAdmin   = "list[admin]"
	InfoStateAdmin   = "info[admin]"
	StopStateAdmin   = "stop[admin]"
	CreateStateAdmin = "create[admin]"
	// PostBackState    = "postback"
)

// type Event string

//	UpgradePremissionEvent = "check premission"
//
// const (
//
//	ListEventUser        = "list instances[user]"
//	ConfigEvent          = "config user info"
//	InfoEventUser        = "show info[user]"
//	StopEventUser        = "stop instances[user]"
//	CreateEventUser      = "create instances[user]"
//	ListEventAdmin       = "list instances[admin]"
//	InfoEventAdmin       = "show info[admin]"
//	StopEventAdmin       = "stop instances[admin]"
//	CreateEventAdmin     = "create instances[admin]"
//	PostBackEvent        = "postback"
//	BackEventUser        = "return to user state[user]"
//	BackEventAdmin       = "return to admin state[admin]"
//
// )
const (
	// UpgradePremissionEvent = "upgrade premission"
	ListEvent   = "list instances"   // -a -n
	ConfigEvent = "config user info" // -u -p -t
	InfoEvent   = "show info"        // -d -n
	StopEvent   = "stop instances"   // -d -n
	CreateEvent = "create instances" // -d -t -n
	// PostBackEvent = "postback"
	BackEvent = "return to original state"
)

func NewUser(id string, con *Controller) *User {
	u := &User{
		UserID:   id,
		UserName: "admin",
		PassWord: "passwd",
		FSM:      stateless.NewStateMachine(UserState),
		Con:      con,
	}
	u.FSM.SetTriggerParameters(ListEvent, reflect.TypeOf(""), reflect.TypeOf(false), reflect.TypeOf(""))
	u.FSM.SetTriggerParameters(ConfigEvent, reflect.TypeOf(""), reflect.TypeOf(""), reflect.TypeOf(""), reflect.TypeOf(""))
	u.FSM.SetTriggerParameters(InfoEvent, reflect.TypeOf(""), reflect.TypeOf(""), reflect.TypeOf(""))
	u.FSM.SetTriggerParameters(StopEvent, reflect.TypeOf(""), reflect.TypeOf(""), reflect.TypeOf(""))
	u.FSM.SetTriggerParameters(CreateEvent, reflect.TypeOf(""), reflect.TypeOf(""), reflect.TypeOf(""), reflect.TypeOf(""))

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(UserState).
		Permit(ListEvent, ListStateUser).
		Permit(ConfigEvent, ConfigState).
		Permit(InfoEvent, InfoStateUser).
		Permit(StopEvent, StopStateUser).
		Permit(CreateEvent, CreateStateUser)
		// Permit(PostBackEvent, PostBackState)

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(AdminState).
		Permit(ListEvent, ListStateUser, func(ctx context.Context, args ...interface{}) bool {
			all := args[1].(bool)
			ns := args[2].(string)
			return canEnterListStateUser(all, ns, u.UserID)
		}).
		Permit(ListEvent, ListStateAdmin, func(ctx context.Context, args ...interface{}) bool {
			all := args[1].(bool)
			ns := args[2].(string)
			return !canEnterListStateUser(all, ns, u.UserID)
		}).
		Permit(ConfigEvent, ConfigState).
		Permit(InfoEvent, InfoStateUser, func(ctx context.Context, args ...interface{}) bool {
			ns := args[2].(string)
			return !canEnterAdmin(ns, u.UserID)
		}).
		Permit(InfoEvent, InfoStateAdmin, func(ctx context.Context, args ...interface{}) bool {
			ns := args[2].(string)
			return canEnterAdmin(ns, u.UserID)
		}).
		Permit(StopEvent, StopStateUser, func(ctx context.Context, args ...interface{}) bool {
			ns := args[2].(string)
			return !canEnterAdmin(ns, u.UserID)
		}).
		Permit(StopEvent, StopStateAdmin, func(ctx context.Context, args ...interface{}) bool {
			ns := args[2].(string)
			return canEnterAdmin(ns, u.UserID)
		}).
		Permit(CreateEvent, CreateStateUser, func(ctx context.Context, args ...interface{}) bool {
			ns := args[3].(string)
			return !canEnterAdmin(ns, u.UserID)

		}).
		Permit(CreateEvent, CreateStateAdmin, func(ctx context.Context, args ...interface{}) bool {
			ns := args[3].(string)
			return canEnterAdmin(ns, u.UserID)
		})
		// Permit(PostBackEvent, PostBackState)

	// ----------------------------------------------------------------------------------------------

	// u.FSM.Configure(PostBackState).
	// 	Permit(BackEvent, UserState, func(ctx context.Context, args ...interface{}) bool {
	// 		return !u.IsAdmin
	// 	}).
	// 	Permit(BackEvent, AdminState, func(ctx context.Context, args ...interface{}) bool {
	// 		return u.IsAdmin
	// 	})

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(ConfigState).
		Permit(BackEvent, UserState, func(ctx context.Context, args ...interface{}) bool {
			return !u.IsAdmin
		}).
		Permit(BackEvent, AdminState, func(ctx context.Context, args ...interface{}) bool {
			return u.IsAdmin
		}).
		OnEntryFrom(ConfigEvent, func(ctx context.Context, args ...interface{}) error {
			replyToken := args[0].(string)
			username := args[1].(string)
			password := args[2].(string)
			admintoken := args[3].(string)
			err := u.handleConfigStateEntry(replyToken, username, password, admintoken)
			if err != nil {
				log.WithField("err", err.Error()).Warn("handleConfigStateEntry")
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(err.Error())).WithContext(ctx).Do()
			} else {
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage("user info configured")).WithContext(ctx).Do()
			}
			if err != nil {
				log.Warn(err.Error())
			}
			return u.FSM.FireCtx(ctx, BackEvent)
		})

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(ListStateUser).
		Permit(BackEvent, UserState, func(ctx context.Context, args ...interface{}) bool {
			return !u.IsAdmin
		}).
		Permit(BackEvent, AdminState, func(ctx context.Context, args ...interface{}) bool {
			return u.IsAdmin
		}).
		OnEntryFrom(ListEvent, func(ctx context.Context, args ...interface{}) error {
			replyToken := args[0].(string)
			list, err := u.handleListStateEntry(ctx, u.UserID)
			if err != nil {
				log.WithField("err", err.Error()).Warn("handleListStateUserEntry")
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(err.Error())).WithContext(ctx).Do()
			} else {
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewFlexMessage("list info", buildListCarousel(list))).WithContext(ctx).Do()
			}
			if err != nil {
				log.Warn(err.Error())
			}
			return u.FSM.FireCtx(ctx, BackEvent)
		})

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(InfoStateUser).
		Permit(BackEvent, UserState, func(ctx context.Context, args ...interface{}) bool {
			return !u.IsAdmin
		}).
		Permit(BackEvent, AdminState, func(ctx context.Context, args ...interface{}) bool {
			return u.IsAdmin
		}).
		OnEntryFrom(InfoEvent, func(ctx context.Context, args ...interface{}) error {
			replyToken := args[0].(string)
			dbname := args[1].(string)
			instance, err := u.handleInfoStateEntry(ctx, u.UserID, dbname)
			if err != nil {
				log.WithField("err", err.Error()).Warn("handleInfoStateUserEntry")
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(err.Error())).WithContext(ctx).Do()
			} else {
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewFlexMessage("info", buildInfoFlexMessage(instance))).WithContext(ctx).Do()
			}
			if err != nil {
				log.Warn(err.Error())
			}
			return u.FSM.FireCtx(ctx, BackEvent)
		})

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(StopStateUser).
		Permit(BackEvent, UserState, func(ctx context.Context, args ...interface{}) bool {
			return !u.IsAdmin
		}).
		Permit(BackEvent, AdminState, func(ctx context.Context, args ...interface{}) bool {
			return u.IsAdmin
		}).
		OnEntryFrom(StopEvent, func(ctx context.Context, args ...interface{}) error {
			replyToken := args[0].(string)
			dbname := args[1].(string)
			// stop target instances
			err := u.handleStopStateEntry(ctx, u.UserID, dbname)
			if err != nil {
				log.WithField("err", err.Error()).Warn("handleStopStateUserEntry")
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(err.Error())).WithContext(ctx).Do()
			} else {
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage("instances stopped")).WithContext(ctx).Do()
			}
			if err != nil {
				log.Warn(err.Error())
			}
			return u.FSM.FireCtx(ctx, BackEvent)
		})

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(CreateStateUser).
		Permit(BackEvent, UserState, func(ctx context.Context, args ...interface{}) bool {
			return !u.IsAdmin
		}).
		Permit(BackEvent, AdminState, func(ctx context.Context, args ...interface{}) bool {
			return u.IsAdmin
		}).
		OnEntryFrom(CreateEvent, func(ctx context.Context, args ...interface{}) error {
			replyToken := args[0].(string)
			dbname := args[1].(string)
			dbtype := args[2].(string)
			err := u.handleCreateStateEntry(ctx, u.UserID, dbtype, dbname)
			if err != nil {
				log.WithField("err", err.Error()).Warn("handleCreateStateUserEntry")
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(err.Error())).WithContext(ctx).Do()
			} else {
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage("instances created")).WithContext(ctx).Do()
			}
			if err != nil {
				log.Warn(err.Error())
			}
			return u.FSM.FireCtx(ctx, BackEvent)
		})

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(ListStateAdmin).
		Permit(BackEvent, AdminState).
		OnEntryFrom(ListEvent, func(ctx context.Context, args ...interface{}) error {
			replyToken := args[0].(string)
			all := args[1].(bool)
			namespace := args[2].(string)
			if all {
				namespace = ""
			} else if namespace == "" {
				namespace = u.UserID
			}
			// list namespace instances
			// get namespace instances
			// reply list of instances or null
			list, err := u.handleListStateEntry(ctx, namespace)
			if err != nil {
				log.WithField("err", err.Error()).Warn("handleListStateAdminEntry")
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(err.Error())).WithContext(ctx).Do()
			} else {
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewFlexMessage("list instances", buildListCarousel(list))).WithContext(ctx).Do()
			}
			if err != nil {
				log.Warn(err.Error())
			}
			return u.FSM.FireCtx(ctx, BackEvent)
		})

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(InfoStateAdmin).
		Permit(BackEvent, AdminState).
		OnEntryFrom(InfoEvent, func(ctx context.Context, args ...interface{}) error {
			// get info of instances dbname
			replyToken := args[0].(string)
			dbname := args[1].(string)
			namespace := args[2].(string)
			if namespace == "" {
				namespace = u.UserID
			}
			instance, err := u.handleInfoStateEntry(ctx, namespace, dbname)
			// get info about dbname and reply
			// info, err := u.Con.k8sclient.GetPodInNamespace(ctx, namespace, dbname)
			if err != nil {
				log.WithField("err", err.Error()).Warn("handleInfoStateAdminEntry")
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(err.Error())).WithContext(ctx).Do()
			} else {
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewFlexMessage("instance info", buildInfoFlexMessage(instance))).WithContext(ctx).Do()
			}
			if err != nil {
				log.Warn(err.Error())
			}
			return u.FSM.FireCtx(ctx, BackEvent)
		})

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(StopStateAdmin).
		Permit(BackEvent, AdminState).
		OnEntryFrom(StopEvent, func(ctx context.Context, args ...interface{}) error {
			replyToken := args[0].(string)
			dbname := args[1].(string)
			namespace := args[2].(string)
			if namespace == "" {
				namespace = u.UserID
			}
			// stop instances dbname
			err := u.handleStopStateEntry(ctx, namespace, dbname)

			// err := u.Con.k8sclient.Delete(ctx, namespace, dbname)
			if err != nil {
				log.WithField("err", err.Error()).Warn("handleStopStateAdminEntry")
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(err.Error())).WithContext(ctx).Do()
			} else {
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage("instances stopped")).WithContext(ctx).Do()
			}
			if err != nil {
				log.Warn(err.Error())
			}

			return u.FSM.FireCtx(ctx, BackEvent)
		})

	// ----------------------------------------------------------------------------------------------

	u.FSM.Configure(CreateStateAdmin).
		Permit(BackEvent, AdminState).
		OnEntryFrom(CreateEvent, func(ctx context.Context, args ...interface{}) error {
			replyToken := args[0].(string)
			dbname := args[1].(string)
			dbtype := args[2].(string)
			namespace := args[3].(string)
			if namespace == "" {
				namespace = u.UserID
			}
			err := u.handleCreateStateEntry(ctx, namespace, dbtype, dbname)
			// create instances dbname
			if err != nil {
				log.WithField("err", err.Error()).Warn("handleCreateStateAdminEntry")
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(err.Error())).WithContext(ctx).Do()
			} else {
				_, err = u.Con.Bot.ReplyMessage(replyToken, linebot.NewTextMessage("instances created")).WithContext(ctx).Do()
			}
			if err != nil {
				log.Warn(err.Error())
			}
			return u.FSM.FireCtx(ctx, BackEvent)
		})

	return u
}

func canEnterListStateUser(all bool, ns string, id string) bool {
	if all {
		return false
	}
	if ns != "" && ns != id {
		return false
	}
	return true
}
func canEnterAdmin(ns string, id string) bool {
	return ns != "" && ns != id
}

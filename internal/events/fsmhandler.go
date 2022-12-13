package events

import (
	"errors"

	"github.com/yungen-lu/TOC-Project-2022/config"
	"github.com/yungen-lu/TOC-Project-2022/internal/models"
	"golang.org/x/net/context"
)

func (u *User) handleConfigStateEntry(replyToken, username, password, admintoken string) error {
	if username != "" {
		u.UserName = username
	}
	if password != "" {
		u.PassWord = password
	}
	if admintoken != "" {
		if admintoken == config.AdminToken {
			u.IsAdmin = true
		} else {
			return errors.New("invalid admin token")
		}
	}
	return nil
}
func (u *User) handleListStateEntry(ctx context.Context, namespace string) ([]models.Instance, error) {
	list, err := u.Con.k8sclient.ListInNamespace(ctx, namespace)
	if err != nil {
		return list, err
	}
	if len(list) == 0 {
		return nil, errors.New("no instances running")
	}
	return list, nil
}
func (u *User) handleInfoStateEntry(ctx context.Context, namespace, dbname string) (*models.Instance, error) {
	// get info about dbname and reply
	info, err := u.Con.k8sclient.GetPodInNamespace(ctx, namespace, dbname)
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (u *User) handleStopStateEntry(ctx context.Context, namespace, dbname string) error {
	return u.Con.k8sclient.Delete(ctx, namespace, dbname)
}
func (u *User) handleCreateStateEntry(ctx context.Context, namespace, dbtype, dbname string) error {
	return u.Con.k8sclient.Create(ctx, namespace, dbtype, dbname, u.UserName, u.PassWord)
}

package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
	"github.com/yungen-lu/TOC-Project-2022/config"
	"github.com/yungen-lu/TOC-Project-2022/internal/controller/http/v1"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

func Run(cfg *config.Config) {
	bot, err := linebot.New(cfg.Line.Secret, cfg.Line.Token)
	if err != nil {
		log.Error(err.Error())
    return
	}
	r := chi.NewRouter()
	v1.NewRouter(r, bot)
	s := http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      r,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}
	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), _defaultShutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}

}

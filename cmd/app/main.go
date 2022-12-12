package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/yungen-lu/TOC-Project-2022/config"
	"github.com/yungen-lu/TOC-Project-2022/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s\n", err.Error())
	}
	app.Run(cfg)
}

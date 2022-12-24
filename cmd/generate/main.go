package main

import (
	"github.com/yungen-lu/TOC-Project-2022/internal/events"
)

func main() {
	u := events.NewUser("testuser", nil)
	print(u.FSM.ToGraph())
}

package main

import (
	// "log"
	// "os"

	// "github.com/looplab/fsm"
	// "github.com/yungen-lu/TOC-Project-2022/internal/state"
	"github.com/yungen-lu/TOC-Project-2022/internal/events"
)

func main() {
	// u := state.NewUser("testuser")
	u := events.NewUser("testuser", nil)
	print(u.FSM.ToGraph())
	// f, err := os.Create("./state_graph.graphviz")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = f.WriteString(fsm.Visualize(u.FSM))
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

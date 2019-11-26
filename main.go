package main

import (
	"os"

	"github.com/gumieri/sway-session/session"
	"github.com/gumieri/typist"
)

var t = typist.New(&typist.Config{})

func main() {
	if len(os.Args) <= 1 {
		return
	}

	switch os.Args[1] {
	case "save":
		s, err := session.New()
		t.Must(err)

		t.Must(s.Save())

	case "restore":
		s, err := session.LoadNewest()
		t.Must(err)

		t.Must(s.Restore())
	}
}

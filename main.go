package main

import (
	"os"
	"strconv"
	"time"

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

	case "save-loop":
		durationInterval := time.Duration(60)
		if len(os.Args) >= 3 {
			interval, err := strconv.Atoi(os.Args[2])
			t.Must(err)
			durationInterval = time.Duration(interval)
		}

		for {
			s, err := session.New()
			t.Must(err)
			time.Sleep(durationInterval * time.Second)
			t.Must(s.Save())
			session.CleanUp()
		}

	case "restore":
		s, err := session.LoadNewest()
		t.Must(err)

		t.Must(s.Restore())

	case "clean-up":
		t.Must(session.CleanUp())
	}
}

package main

import (
	"encoding/gob"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/gumieri/go-sway"
	"github.com/gumieri/typist"
)

var t = typist.New(&typist.Config{})

var sessionFilePath = path.Join(xdg.DataHome, "sway-session", "session.gob")

func saveSession() (err error) {
	tree, err := sway.GetTree()
	if err != nil {
		return
	}

	programs, err := GetPrograms(&GetProgramsInput{
		Parent: tree.Root,
		Procs:  AllProcs(),
	})
	if err != nil {
		return
	}

	sessionFile, err := os.Create(sessionFilePath)
	if err != nil {
		return
	}

	dataEncoder := gob.NewEncoder(sessionFile)
	dataEncoder.Encode(programs)

	sessionFile.Close()

	return
}

func loadSession() (programs []*Program, err error) {
	sessionFile, err := os.Open(sessionFilePath)

	if err != nil {
		return
	}

	dataDecoder := gob.NewDecoder(sessionFile)
	err = dataDecoder.Decode(&programs)
	if err != nil {
		return
	}

	sessionFile.Close()

	return
}

func main() {
	if len(os.Args) <= 1 {
		return
	}

	switch os.Args[1] {
	case "save":
		t.Must(saveSession())
	case "restore":
		programs, err := loadSession()
		t.Must(err)

		for _, program := range programs {
			_, err := sway.RunCommand(program.Restore())
			t.Must(err)
		}
	}
}

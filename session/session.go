package session

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/adrg/xdg"
	"github.com/gumieri/go-sway"
	"github.com/gumieri/sway-session/proc"
	"github.com/gumieri/sway-session/program"
)

var configPath = path.Join(xdg.DataHome, "sway-session")
var sessionsPath = path.Join(configPath, "sessions")

// Session have information about a Sway Session
type Session struct {
	FilePath string             `json:"string"`
	Programs []*program.Program `json:"programs"`
}

// New create a instance of Session
func New() (s *Session, err error) {
	filename := strconv.FormatInt(time.Now().Unix(), 10) + ".json"
	filePath := path.Join(sessionsPath, filename)
	s = &Session{FilePath: filePath}

	tree, err := sway.GetTree()
	if err != nil {
		return
	}

	s.Programs, err = program.GetPrograms(&program.GetProgramsInput{
		Parent: tree.Root,
		Procs:  proc.AllProcs(),
	})

	return
}

func timestampFromFilename(filename string) int {
	s := filename[0 : len(filename)-len(filepath.Ext(filename))]
	t, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return t
}

// LoadNewest read the newest saved Session data from the disk
func LoadNewest() (s *Session, err error) {
	files, err := ioutil.ReadDir(sessionsPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		err = errors.New("No session file found")
		return
	}

	sort.Slice(files, func(i, j int) bool {
		iUnix := timestampFromFilename(files[i].Name())
		jUnix := timestampFromFilename(files[j].Name())
		return iUnix > jUnix
	})

	filePath := path.Join(sessionsPath, files[0].Name())
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &s)

	return
}

// Save write to disk the Session data
func (s *Session) Save() (err error) {
	b, err := json.Marshal(s)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(s.FilePath, b, 0600)

	return
}

// Restore execute the commands at Sway to restore its state
func (s *Session) Restore() (err error) {
	for _, p := range s.Programs {
		_, err = sway.RunCommand(p.Restore())
		if err != nil {
			return
		}
	}

	return
}

// CleanUp delete all session files except the last
func CleanUp() (err error) {
	files, err := ioutil.ReadDir(sessionsPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		return
	}

	sort.Slice(files, func(i, j int) bool {
		iUnix := timestampFromFilename(files[i].Name())
		jUnix := timestampFromFilename(files[j].Name())
		return iUnix > jUnix
	})

	for i, file := range files {
		if i == 0 {
			continue
		}

		err = os.Remove(path.Join(sessionsPath, file.Name()))
		if err != nil {
			return
		}
	}

	return
}

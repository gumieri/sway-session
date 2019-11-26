package proc

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Procs is a group of processes
type Procs []*Proc

// Proc maps the data about a system process
type Proc struct {
	PID     int
	PPID    int
	EXE     string
	CWD     string
	CMDLine []string
	Stat    []string
}

// NewProc gather the system process data from a informed PID
func NewProc(pid int) (p *Proc, err error) {
	p = &Proc{PID: pid}

	pidS := strconv.Itoa(pid)
	procDir := "/proc/" + pidS + "/"

	p.EXE, err = os.Readlink(procDir + "exe")
	if err != nil {
		return
	}

	p.CWD, err = os.Readlink(procDir + "cwd")
	if err != nil {
		return
	}

	cmdline, err := ioutil.ReadFile(procDir + "cmdline")
	if err != nil {
		return
	}
	p.CMDLine = strings.Split(string(cmdline), "\x00")

	stat, err := ioutil.ReadFile(procDir + "stat")
	if err != nil {
		return
	}
	p.Stat = strings.Split(string(stat), " ")

	p.PPID, _ = strconv.Atoi(p.Stat[3])

	return
}

// AllProcs gather all running procces at the system to return a collection of it
func AllProcs() *Procs {
	files, err := ioutil.ReadDir("/proc/")
	if err != nil {
		return nil
	}

	var procs Procs
	for _, f := range files {
		pid, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}

		proc, err := NewProc(pid)
		if err != nil {
			continue
		}

		procs = append(procs, proc)
	}

	return &procs
}

// Find a informed PID in the collection of Procs
func (ps *Procs) Find(pid int) *Proc {
	for _, p := range *ps {
		if p.PID == pid {
			return p
		}
	}

	return nil
}

// ChildrenOf look for all process having the informed Proc as parent
func (ps *Procs) ChildrenOf(pp *Proc) *Procs {
	var c Procs
	for _, p := range *ps {
		if p.PPID == pp.PID {
			c = append(c, p)
		}
	}

	return &c
}

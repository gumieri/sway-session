package program

import (
	"errors"
	"strings"

	"github.com/gumieri/go-sway"
	"github.com/gumieri/sway-session/proc"
)

// Program define the information to recreate a program
type Program struct {
	Name      string   `json:"name"`
	Workspace string   `json:"workspace"`
	Command   []string `json:"command"`
}

// NewProgramInput is only used as parameter data to NewProgram
type NewProgramInput struct {
	Node      *sway.Node
	Workspace string
	Procs     *proc.Procs
}

// NewProgram create a Program with the definitions of informed Sway Node, for a defined Workspace
func NewProgram(input *NewProgramInput) (p *Program, err error) {
	p = &Program{Workspace: input.Workspace}

	proc := input.Procs.Find(input.Node.PID)
	if proc == nil {
		err = errors.New("PID not found")
		return
	}

	exeA := strings.Split(proc.EXE, "/")
	p.Name = exeA[len(exeA)-1]

	p.Command = proc.CMDLine
	switch p.Name {
	case "alacritty":
		children := *input.Procs.ChildrenOf(proc)
		p.Command = []string{proc.CMDLine[0], "--working-directory " + children[0].CWD}
	case "firefox":
		p.Command = []string{"/usr/bin/firefox"}
	}

	return
}

// Restore outputs the command for restoring a Program
func (p *Program) Restore() string {
	return "workspace " + p.Workspace + "; exec " + strings.Join(p.Command, " ")
}

// GetProgramsInput is only used as parameter data to GetPrograms
type GetProgramsInput struct {
	Parent    *sway.Node
	Workspace string
	Procs     *proc.Procs
}

// GetPrograms read a Sway Tree for mapping the running programs and return a slice of it
func GetPrograms(input *GetProgramsInput) ([]*Program, error) {
	programs := make([]*Program, 0)

	for _, node := range input.Parent.Nodes {
		switch node.Type {
		case sway.Con:
			p, err := NewProgram(&NewProgramInput{
				Node:      node,
				Workspace: input.Workspace,
				Procs:     input.Procs,
			})

			if err != nil {
				return programs, err
			}

			programs = append(programs, p)

		case sway.WorkspaceNode:
			input.Workspace = node.Name
			fallthrough

		default:
			nodePrograms, err := GetPrograms(&GetProgramsInput{
				Parent:    node,
				Workspace: input.Workspace,
				Procs:     input.Procs,
			})

			if err != nil {
				return programs, err
			}

			programs = append(programs, nodePrograms...)
		}
	}

	return programs, nil
}

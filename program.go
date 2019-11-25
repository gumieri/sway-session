package main

import (
	"errors"
	"strings"

	"github.com/gumieri/go-sway"
)

type Program struct {
	Proc      *Proc
	Name      string
	Workspace string
	Command   []string
}

type NewProgramInput struct {
	Node      *sway.Node
	Workspace string
	Procs     *Procs
}

func NewProgram(input *NewProgramInput) (p *Program, err error) {
	p = &Program{Workspace: input.Workspace}

	p.Proc = input.Procs.Find(input.Node.PID)
	if p.Proc == nil {
		err = errors.New("PID not found")
		return
	}

	exeA := strings.Split(p.Proc.EXE, "/")
	p.Name = exeA[len(exeA)-1]

	p.Command = p.Proc.CMDLine
	switch p.Name {
	case "alacritty":
		children := *input.Procs.ChildrenOf(p.Proc)
		p.Command = []string{
			p.Proc.CMDLine[0],
			"--working-directory " + children[0].CWD,
		}
	}

	return
}

func (p *Program) Restore() string {
	return "workspace " + p.Workspace + "; exec " + strings.Join(p.Command, " ")
}

type GetProgramsInput struct {
	Parent    *sway.Node
	Workspace string
	Procs     *Procs
}

func GetPrograms(input *GetProgramsInput) (programs []*Program, err error) {
	programs = make([]*Program, 0)

	for _, node := range input.Parent.Nodes {
		switch node.Type {
		case sway.Con:
			var p *Program
			p, err = NewProgram(&NewProgramInput{
				Node:      node,
				Workspace: input.Workspace,
				Procs:     input.Procs,
			})

			if err != nil {
				return
			}

			programs = append(programs, p)

		case sway.WorkspaceNode:
			input.Workspace = node.Name
			fallthrough

		default:
			nodePrograms := make([]*Program, 0)
			nodePrograms, err = GetPrograms(&GetProgramsInput{
				Parent:    node,
				Workspace: input.Workspace,
				Procs:     input.Procs,
			})

			if err != nil {
				return
			}

			programs = append(programs, nodePrograms...)
		}
	}

	return
}

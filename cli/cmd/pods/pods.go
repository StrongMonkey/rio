package pods

import (
	"github.com/rancher/rio/cli/cmd/ps"
	"github.com/rancher/rio/cli/pkg/builder"
	"github.com/rancher/rio/cli/pkg/clicontext"
	"github.com/rancher/rio/cli/pkg/table"
	"github.com/urfave/cli"
)

type Pod struct {
	C_Containers bool `desc:"print containers, not services"`
}

func (p *Pod) Customize(cmd *cli.Command) {
	cmd.Flags = append(table.WriterFlags(), cmd.Flags...)
}

func Pods(app *cli.App) cli.Command {
	ls := builder.Command(&Pod{},
		"List Pods",
		app.Name+" pods [OPTIONS]",
		"")
	return cli.Command{
		Name:      "pods",
		ShortName: "pod",
		Usage:     "Operations on Pods",
		Action:    clicontext.DefaultAction(ls.Action),
		Flags:     ls.Flags,
		Category:  "SUB COMMANDS",
		Subcommands: []cli.Command{
			ls,
		},
	}
}

func (p *Pod) Run(ctx *clicontext.CLIContext) error {
	if p.C_Containers {
		return ps.Containers(ctx)
	}
	return ps.Pods(ctx)
}

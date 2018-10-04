package exec

import (
	"fmt"

	"github.com/rancher/rio/cli/cmd/ps"
	"github.com/rancher/rio/cli/pkg/clicontext"
	"github.com/sirupsen/logrus"
)

type Exec struct {
	I_Stdin     bool   `desc:"Pass stdin to the container"`
	T_Tty       bool   `desc:"Stdin is a TTY"`
	C_Container string `desc:"Specific container in pod, default is first container"`
}

func (e *Exec) Run(ctx *clicontext.CLIContext) error {
	cc, err := ctx.ClusterClient()
	if err != nil {
		return err
	}

	args := ctx.CLI.Args()
	if len(args) < 2 {
		return fmt.Errorf("at least two arguments are required CONTAINER CMD")
	}

	cd, err := ps.ListFirstPod(cc, true, e.C_Container, args[0])
	if err != nil {
		return err
	}

	if cd == nil {
		return fmt.Errorf("failed to find pod for %s, container \"%s\"", args[0], e.C_Container)
	}

	podNS, podName, containerName := cd.Pod.Namespace, cd.Pod.Name, cd.Container.Name

	execArgs := []string{}
	if logrus.GetLevel() >= logrus.DebugLevel {
		execArgs = append(execArgs, "-v=9")
	}
	if e.I_Stdin {
		execArgs = append(execArgs, "-i")
	}
	if e.T_Tty {
		execArgs = append(execArgs, "-t")
	}

	execArgs = append(execArgs, podName, "-c", containerName)
	execArgs = append(execArgs, args[1:]...)

	cluster, err := ctx.Cluster()
	if err != nil {
		return err
	}
	return cluster.Kubectl(podNS, "exec", execArgs...)
}

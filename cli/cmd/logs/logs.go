package logs

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/rancher/rio/cli/cmd/ps"
	"github.com/rancher/rio/cli/pkg/clicontext"
	"github.com/rancher/rio/cli/pkg/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Logs struct {
	F_Follow     bool   `desc:"Follow log output"`
	S_Since      string `desc:"Logs since a certain time, either duration (5s, 2m, 3h) or RFC3339"`
	P_Previous   bool   `desc:"Print the logs for the previous instance of the container in a pod if it exists"`
	C_Container  string `desc:"Print the logs of a specific container"`
	R_Revision   string `desc:"Print the logs of a specific revision"`
	N_Tail       int    `desc:"Number of recent lines of logs to print, -1 for all" default:"200"`
	A_All        bool   `desc:"Include hidden or systems logs when logging"`
	T_Timestamps bool   `desc:"Print the logs with timestamp"`
}

func (l *Logs) Run(ctx *clicontext.CLIContext) error {
	if len(ctx.CLI.Args()) == 0 {
		return fmt.Errorf("at least one argument is required: CONTAINER_OR_SERVICE")
	}

	logPods := false
	for _, arg := range ctx.CLI.Args() {
		if strings.Count(arg, "/") == 2 {
			logPods = true
		}
	}

	pds, err := ps.ListPods(ctx, true, ctx.CLI.Args()...)
	if err != nil {
		return err
	}

	if len(pds) == 0 {
		return fmt.Errorf("failed to find container for %v, container \"%s\"", ctx.CLI.Args(), l.C_Container)
	}

	factory := logger.NewColorLoggerFactory()
	errg, _ := errgroup.WithContext(ctx.Ctx)
	for _, pd := range pds {
		if l.R_Revision != "" && pd.Service.Version != l.R_Revision {
			continue
		}
		for _, container := range pd.Containers {
			pod := pd.Pod
			c := container
			if l.C_Container != "" && container.Name != l.C_Container {
				continue
			}
			if !l.A_All && (container.Name == "linkerd-proxy" || container.Name == "linkerd-init") {
				if l.C_Container == "" && !logPods {
					continue
				}
			}
			errg.Go(func() error {
				return l.logContainer(pod, c, ctx.Core, factory)
			})
		}
	}
	return errg.Wait()
}

func (l *Logs) logContainer(pod *v1.Pod, container v1.Container, coreClient corev1.CoreV1Interface, factory *logger.ColorLoggerFactory) error {
	containerName := fmt.Sprintf("%s/%s", pod.Name, container.Name)
	logger := factory.CreateContainerLogger(containerName)
	podLogOption := &v1.PodLogOptions{
		Container: container.Name,
		Follow:    l.F_Follow,
	}
	if l.T_Timestamps {
		podLogOption.Timestamps = true
	}
	if l.S_Since != "" {
		t, err := time.Parse(time.RFC3339, l.S_Since)
		if err == nil {
			newTime := metav1.NewTime(t)
			podLogOption.SinceTime = &newTime
		} else {
			du, err := time.ParseDuration(l.S_Since)
			if err == nil {
				ss := int64(du.Round(time.Second).Seconds())
				podLogOption.SinceSeconds = &ss
			}
		}
	}

	req := coreClient.Pods(pod.Namespace).GetLogs(pod.Name, podLogOption)
	reader, err := req.Stream()
	if err != nil {
		return err
	}
	defer reader.Close()

	sc := bufio.NewScanner(reader)
	for sc.Scan() {
		logger.Out(append(sc.Bytes(), []byte("\n")...))
	}

	return nil
}

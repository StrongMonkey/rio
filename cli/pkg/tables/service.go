package tables

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rancher/rio/cli/pkg/table"
	riov1 "github.com/rancher/rio/pkg/apis/rio.cattle.io/v1"
	v1 "github.com/rancher/rio/pkg/apis/rio.cattle.io/v1"
	"github.com/rancher/rio/pkg/services"
	corev1 "k8s.io/api/core/v1"
)

func NewService(cfg Config) TableWriter {
	writer := table.NewWriter([][]string{
		{"NAME", "{{.ID}}"},
		{"IMAGE", "{{.Service | image}}"},
		{"ENDPOINT", "{{arrayFirst .Service.Status.Endpoints}}"},
		{"SCALE", "{{scale .Service .Service.Status.ScaleStatus}}"},
		{"APP", "{{.Service | app}}"},
		{"VERSION", "{{.Service | version}}"},
		{"WEIGHT", "{{.Service | formatWeight}}"},
		{"CREATED", "{{.Service.CreationTimestamp | ago}}"},
		{"DETAIL", "{{serviceDetail .Service .Pod}}"},
	}, cfg)

	writer.AddFormatFunc("image", FormatImage)
	writer.AddFormatFunc("scale", formatRevisionScale)
	writer.AddFormatFunc("app", app)
	writer.AddFormatFunc("version", version)
	writer.AddFormatFunc("formatWeight", formatWeight)
	writer.AddFormatFunc("serviceDetail", serviceDetail)

	return &tableWriter{
		writer: writer,
	}
}

func app(data interface{}) string {
	s, ok := data.(*v1.Service)
	if !ok {
		return ""
	}
	appName, _ := services.AppAndVersion(s)
	return appName
}

func version(data interface{}) string {
	s, ok := data.(*v1.Service)
	if !ok {
		return ""
	}
	_, version := services.AppAndVersion(s)
	return version
}

func formatWeight(data interface{}) string {
	s, ok := data.(*v1.Service)
	if !ok {
		return ""
	}

	if s.Status.ComputedWeight != nil {
		return fmt.Sprintf("%s%%", strconv.Itoa(*s.Status.ComputedWeight))
	}
	return "0%"
}

func serviceDetail(data interface{}, pod *corev1.Pod) string {
	s, ok := data.(*v1.Service)
	if !ok {
		return ""
	}

	for _, con := range s.Status.Conditions {
		if con.Status != corev1.ConditionTrue {
			return fmt.Sprintf("%s: %s(%s)", con.Type, con.Message, con.Reason)
		}
	}

	output := strings.Builder{}
	if pod == nil {
		return ""
	}
	for _, con := range append(pod.Status.ContainerStatuses, pod.Status.InitContainerStatuses...) {
		if con.State.Waiting != nil && con.State.Waiting.Reason != "" {
			output.WriteString("; ")
			reason := con.State.Waiting.Reason
			if con.State.Waiting.Message != "" {
				reason = reason + "/" + con.State.Waiting.Message
			}
			output.WriteString(fmt.Sprintf("%s(%s)", con.Name, reason))
		}

		if con.State.Terminated != nil && con.State.Terminated.ExitCode != 0 {
			output.WriteString(";")
			if con.State.Terminated.Message == "" {
				con.State.Terminated.Message = "exit code not zero"
			}
			reason := con.State.Terminated.Reason
			if con.State.Terminated.Message != "" {
				reason = reason + "/" + con.State.Terminated.Message
			}
			output.WriteString(fmt.Sprintf("%s(%s), exit code: %v", con.Name, reason, con.State.Terminated.ExitCode))
		}
	}
	return strings.Trim(output.String(), "; ")
}

func formatRevisionScale(svc *riov1.Service, scaleStatus *v1.ScaleStatus) (string, error) {
	scale := svc.Spec.Replicas
	if svc.Status.ComputedReplicas != nil && services.AutoscaleEnable(svc) {
		scale = svc.Status.ComputedReplicas
	}
	return FormatScale(scale, scaleStatus)
}

func FormatScale(scale *int, scaleStatus *v1.ScaleStatus) (string, error) {
	scaleNum := 1
	if scale != nil {
		scaleNum = *scale
	}

	scaleStr := strconv.Itoa(scaleNum)

	if scaleStatus == nil {
		scaleStatus = &v1.ScaleStatus{}
	}

	if scaleNum == -1 {
		return strconv.Itoa(scaleStatus.Available), nil
	}

	if scaleStatus.Unavailable == 0 {
		return scaleStr, nil
	}

	var prefix string
	percentage := ""
	ready := scaleNum - scaleStatus.Unavailable
	if scaleNum > 0 {
		percentage = fmt.Sprintf(" %d%%", (ready*100)/scaleNum)
	}

	if ready != scaleNum {
		prefix = fmt.Sprintf("%d/", ready)
	}

	return fmt.Sprintf("%s%d%s", prefix, scaleNum, percentage), nil
}

func FormatImage(data interface{}) (string, error) {
	s, ok := data.(*v1.Service)
	if !ok {
		return fmt.Sprint(data), nil
	}
	image := ""
	if s.Spec.Image == "" && len(s.Spec.Sidecars) > 0 {
		image = s.Spec.Sidecars[0].Image
	} else {
		image = s.Spec.Image
	}
	return strings.TrimPrefix(image, "localhost:5442/"), nil
}

package tables

import (
	"fmt"

	"github.com/rancher/rio/cli/pkg/table"
	adminv1 "github.com/rancher/rio/pkg/apis/admin.rio.cattle.io/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewPublicDomain(cfg Config) TableWriter {
	writer := table.NewWriter([][]string{
		{"NAME", "{{id .Obj}}"},
		{"TARGET", "{{targetName .Obj}}"},
	}, cfg)

	writer.AddFormatFunc("targetName", formatTargetName)

	return &tableWriter{
		writer: writer,
	}
}

func formatTargetName(obj runtime.Object) string {
	domain := obj.(*adminv1.PublicDomain)

	writeNs := true
	if domain.Namespace == domain.Spec.TargetNamespace {
		writeNs = false
	}

	target := domain.Spec.TargetApp
	if domain.Spec.TargetVersion != "" {
		target = fmt.Sprintf("%s:%s", target, domain.Spec.TargetVersion)
	}

	if writeNs {
		target = fmt.Sprintf("%s/%s", domain.Spec.TargetNamespace, target)
	}
	return target
}

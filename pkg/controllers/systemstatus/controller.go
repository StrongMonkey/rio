package systemstatus

import (
	"context"

	riov1 "github.com/rancher/rio/pkg/apis/rio.cattle.io/v1"
	adminv1controller "github.com/rancher/rio/pkg/generated/controllers/admin.rio.cattle.io/v1"
	"github.com/rancher/rio/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Register(ctx context.Context, rContext *types.Context) error {
	h := handler{
		systemNamespace: rContext.Namespace,
		rioinfos:        rContext.Admin.Admin().V1().RioInfo(),
	}

	rContext.Rio.Rio().V1().Service().OnChange(ctx, "systemstatus-check", h.sync)
	return nil
}

type handler struct {
	systemNamespace string
	rioinfos        adminv1controller.RioInfoController
}

func (h handler) sync(key string, obj *riov1.Service) (*riov1.Service, error) {
	if obj == nil || obj.DeletionTimestamp != nil {
		return obj, nil
	}

	if obj.Namespace != h.systemNamespace {
		return obj, nil
	}

	infoObj, err := h.rioinfos.Get("rio", metav1.GetOptions{})
	if err != nil {
		return obj, err
	}

	info := infoObj.DeepCopy()
	readyMap := info.Status.SystemComponentReadyMap
	if readyMap == nil {
		readyMap = make(map[string]string)
	}
	if obj.Status.ScaleStatus != nil && obj.Status.ScaleStatus.Available > 0 {
		readyMap[obj.Name] = "running"
	}
	info.Status.SystemComponentReadyMap = readyMap
	if _, err := h.rioinfos.Update(info); err != nil {
		return obj, err
	}
	return obj, nil
}

package istio

import (
	"context"

	"github.com/rancher/norman/pkg/changeset"
	"github.com/rancher/norman/pkg/objectset"
	"github.com/rancher/rio/features/letsencrypt/controllers/issuer"
	"github.com/rancher/rio/features/routing/controllers/istio/populate"
	"github.com/rancher/rio/pkg/settings"
	"github.com/rancher/rio/types"
	"github.com/rancher/rio/types/apis/networking.istio.io/v1alpha3"
	projectv1 "github.com/rancher/rio/types/apis/project.rio.cattle.io/v1"
	corev1 "github.com/rancher/types/apis/core/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	all = "_istio_deploy_"
)

var trigger = []changeset.Key{
	{
		Name: all,
	},
}

func Register(ctx context.Context, rContext *types.Context) error {
	s := &istioDeployController{
		template: objectset.NewProcessor("istio-gateway").
			Client(rContext.Rio.Stack,
				rContext.Networking.Gateway),
		publicdomainLister: rContext.Global.PublicDomain.Cache(),
		secretsLister:      rContext.Core.Secret.Cache(),
	}

	rContext.Networking.VirtualService.Interface().AddHandler(ctx, "istio-deploy", s.sync)
	changeset.Watch(ctx, "istio-deploy",
		resolve,
		rContext.Networking.VirtualService,
		rContext.Networking.VirtualService,
		rContext.Core.Namespace)
	rContext.Networking.VirtualService.Enqueue("", all)

	return nil
}

func resolve(namespace, name string, obj runtime.Object) ([]changeset.Key, error) {
	switch t := obj.(type) {
	case *v1alpha3.VirtualService:
		return trigger, nil
	case *v1.Namespace:
		if t.Name == settings.IstioExternalLBNamespace {
			return trigger, nil
		}
	}

	return nil, nil
}

type istioDeployController struct {
	template           objectset.Processor
	publicdomainLister projectv1.PublicDomainClientCache
	secretsLister      corev1.SecretClientCache
}

func (i *istioDeployController) sync(key string, obj *v1alpha3.VirtualService) (runtime.Object, error) {
	if key != all {
		return nil, nil
	}

	pds, err := i.publicdomainLister.List("", labels.Everything())
	if err != nil {
		return nil, err
	}

	secret, err := i.secretsLister.Get(settings.IstioExternalLBNamespace, issuer.TLSSecretName)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	os := populate.Istio(pds, secret)
	return nil, i.template.NewDesiredSet(nil, os).Apply()
}

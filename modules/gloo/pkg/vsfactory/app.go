package vsfactory

import (
	riov1 "github.com/rancher/rio/pkg/apis/rio.cattle.io/v1"
	solov1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
)

func (f *VirtualServiceFactory) ForApp(namespace, appName string, svcs []*riov1.Service) ([]*solov1.VirtualService, error) {
	hostnames, targets, err := getTargetsForApp(svcs)
	if err != nil {
		return nil, err
	}

	vs := newVirtualService(namespace, appName, hostnames)
	vs.Spec.VirtualHost.Routes[0].Action = newRouteAction(targets...)

	result := []*solov1.VirtualService{
		vs,
	}

	tls, err := f.findTLS(namespace, appName, "", hostnames)
	if err != nil {
		return nil, err
	}

	for hostname, tls := range tls {
		result = append(result, tlsCopy(hostname, f.systemNamespace, tls, vs))
	}

	return result, nil
}

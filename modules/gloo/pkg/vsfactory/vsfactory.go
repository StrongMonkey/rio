package vsfactory

import (
	rioadminv1controller "github.com/rancher/rio/pkg/generated/controllers/admin.rio.cattle.io/v1"
	"github.com/rancher/rio/types"
	"github.com/rancher/wrangler/pkg/name"
	soloapiv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solov1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualServiceFactory struct {
	clusterDomainCache rioadminv1controller.ClusterDomainCache
	publicDomainCache  rioadminv1controller.PublicDomainCache
	systemNamespace    string
}

func New(rContext *types.Context) *VirtualServiceFactory {
	return &VirtualServiceFactory{
		clusterDomainCache: rContext.Admin.Admin().V1().ClusterDomain().Cache(),
		publicDomainCache:  rContext.Admin.Admin().V1().PublicDomain().Cache(),
		systemNamespace:    rContext.Namespace,
	}
}

func newVirtualService(namespace, name string, hosts []string) *solov1.VirtualService {
	return &solov1.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: soloapiv1.VirtualService{
			VirtualHost: &soloapiv1.VirtualHost{
				Domains: hosts,
				Routes: []*soloapiv1.Route{
					{
						Matcher: &gloov1.Matcher{
							PathSpecifier: &gloov1.Matcher_Prefix{
								Prefix: "/",
							},
						},
					},
				},
			},
		},
	}
}

func tlsCopy(hostname, tlsNamespace, tlsName string, vs *solov1.VirtualService) *solov1.VirtualService {
	vsTLS := vs.DeepCopy()
	vsTLS.Name = name.SafeConcatName(vsTLS.Name, "tls", tlsName)
	vsTLS.Spec.SslConfig = &gloov1.SslConfig{
		SslSecrets: &gloov1.SslConfig_SecretRef{
			SecretRef: &core.ResourceRef{
				Name:      tlsName,
				Namespace: tlsNamespace,
			},
		},
	}
	vsTLS.Spec.VirtualHost.Domains = []string{hostname}
	return vsTLS
}

func newRouteAction(targets ...target) *soloapiv1.Route_RouteAction {
	if len(targets) == 0 {
		return nil
	}

	if len(targets) == 1 {
		return single(targets[0])
	}

	var dest []*gloov1.WeightedDestination

	for _, target := range targets {
		dest = append(dest, &gloov1.WeightedDestination{
			Destination: destination(target),
			Weight:      uint32(target.Weight),
		})
	}

	return &soloapiv1.Route_RouteAction{
		RouteAction: &gloov1.RouteAction{
			Destination: &gloov1.RouteAction_Multi{
				Multi: &gloov1.MultiDestination{
					Destinations: dest,
				},
			},
		},
	}
}

func destination(target target) *gloov1.Destination {
	return &gloov1.Destination{
		DestinationType: &gloov1.Destination_Kube{
			Kube: &gloov1.KubernetesServiceDestination{
				Ref: core.ResourceRef{
					Name:      target.Name,
					Namespace: target.Namespace,
				},
				Port: uint32(target.Port),
			},
		},
	}
}

func single(target target) *soloapiv1.Route_RouteAction {
	return &soloapiv1.Route_RouteAction{
		RouteAction: &gloov1.RouteAction{
			Destination: &gloov1.RouteAction_Single{
				Single: destination(target),
			},
		},
	}
}

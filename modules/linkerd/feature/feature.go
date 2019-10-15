package feature

import (
	"context"

	"github.com/rancher/rio/modules/linkerd/controllers/app"
	"github.com/rancher/rio/modules/linkerd/controllers/router"
	"github.com/rancher/rio/modules/linkerd/pkg/injector"
	"github.com/rancher/rio/pkg/features"
	"github.com/rancher/rio/pkg/stack"
	"github.com/rancher/rio/types"
)

func Register(ctx context.Context, rContext *types.Context) error {
	apply := rContext.Apply.WithCacheTypes(rContext.Rio.Rio().V1().Router())
	feature := &features.FeatureController{
		FeatureName: "linkerd",
		FeatureSpec: features.FeatureSpec{
			Description: "linkerd service mesh",
		},
		SystemStacks: []*stack.SystemStack{
			stack.NewSystemStack(apply, "linkerd", "linkerd"),
		},
		Controllers: []features.ControllerRegister{
			app.Register,
			router.Register,
		},
		OnStart: func() error {
			injector.RegisterInjector()
			rContext.Rio.Rio().V1().Service().Enqueue("*", "*")
			return nil
		},
	}
	return feature.Register()
}

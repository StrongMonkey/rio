package modules

import (
	"context"

	"github.com/rancher/rio/modules/smi"

	"github.com/rancher/rio/pkg/indexes"

	"github.com/rancher/rio/modules/gloo"
	"github.com/rancher/rio/modules/info"
	"github.com/rancher/rio/modules/letsencrypt"
	"github.com/rancher/rio/modules/rdns"
	"github.com/rancher/rio/modules/service"
	"github.com/rancher/rio/types"
)

func Register(ctx context.Context, rContext *types.Context) error {
	indexes.RegisterIndexes(rContext)

	if err := info.Register(ctx, rContext); err != nil {
		return err
	}
	if err := rdns.Register(ctx, rContext); err != nil {
		return err
	}
	if err := service.Register(ctx, rContext); err != nil {
		return err
	}
	if err := gloo.Register(ctx, rContext); err != nil {
		return err
	}
	if err := smi.Register(ctx, rContext); err != nil {
		return err
	}
	//if err := linkerd.Register(ctx, rioContext); err != nil {
	//	return err
	//}
	if err := letsencrypt.Register(ctx, rContext); err != nil {
		return err
	}
	//if err := build.Register(ctx, rioContext); err != nil {
	//	return err
	//}
	//if err := autoscale.Register(ctx, rioContext); err != nil {
	//	return err
	//}
	return nil
}

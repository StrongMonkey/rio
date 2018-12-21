package v1

import (
	"context"
	"sync"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"github.com/rancher/norman/objectclient/dynamic"
	"github.com/rancher/norman/restwatch"
	"k8s.io/client-go/rest"
)

type (
	contextKeyType        struct{}
	contextClientsKeyType struct{}
)

type Interface interface {
	RESTClient() rest.Interface
	controller.Starter

	ListenConfigsGetter
	SettingsGetter
	PublicDomainsGetter
	FeaturesGetter
}

type Clients struct {
	Interface Interface

	ListenConfig ListenConfigClient
	Setting      SettingClient
	PublicDomain PublicDomainClient
	Feature      FeatureClient
}

type Client struct {
	sync.Mutex
	restClient rest.Interface
	starters   []controller.Starter

	listenConfigControllers map[string]ListenConfigController
	settingControllers      map[string]SettingController
	publicDomainControllers map[string]PublicDomainController
	featureControllers      map[string]FeatureController
}

func Factory(ctx context.Context, config rest.Config) (context.Context, controller.Starter, error) {
	c, err := NewForConfig(config)
	if err != nil {
		return ctx, nil, err
	}

	cs := NewClientsFromInterface(c)

	ctx = context.WithValue(ctx, contextKeyType{}, c)
	ctx = context.WithValue(ctx, contextClientsKeyType{}, cs)
	return ctx, c, nil
}

func ClientsFrom(ctx context.Context) *Clients {
	return ctx.Value(contextClientsKeyType{}).(*Clients)
}

func From(ctx context.Context) Interface {
	return ctx.Value(contextKeyType{}).(Interface)
}

func NewClients(config rest.Config) (*Clients, error) {
	iface, err := NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return NewClientsFromInterface(iface), nil
}

func NewClientsFromInterface(iface Interface) *Clients {
	return &Clients{
		Interface: iface,

		ListenConfig: &listenConfigClient2{
			iface: iface.ListenConfigs(""),
		},
		Setting: &settingClient2{
			iface: iface.Settings(""),
		},
		PublicDomain: &publicDomainClient2{
			iface: iface.PublicDomains(""),
		},
		Feature: &featureClient2{
			iface: iface.Features(""),
		},
	}
}

func NewForConfig(config rest.Config) (Interface, error) {
	if config.NegotiatedSerializer == nil {
		config.NegotiatedSerializer = dynamic.NegotiatedSerializer
	}

	restClient, err := restwatch.UnversionedRESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{
		restClient: restClient,

		listenConfigControllers: map[string]ListenConfigController{},
		settingControllers:      map[string]SettingController{},
		publicDomainControllers: map[string]PublicDomainController{},
		featureControllers:      map[string]FeatureController{},
	}, nil
}

func (c *Client) RESTClient() rest.Interface {
	return c.restClient
}

func (c *Client) Sync(ctx context.Context) error {
	return controller.Sync(ctx, c.starters...)
}

func (c *Client) Start(ctx context.Context, threadiness int) error {
	return controller.Start(ctx, threadiness, c.starters...)
}

type ListenConfigsGetter interface {
	ListenConfigs(namespace string) ListenConfigInterface
}

func (c *Client) ListenConfigs(namespace string) ListenConfigInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &ListenConfigResource, ListenConfigGroupVersionKind, listenConfigFactory{})
	return &listenConfigClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type SettingsGetter interface {
	Settings(namespace string) SettingInterface
}

func (c *Client) Settings(namespace string) SettingInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &SettingResource, SettingGroupVersionKind, settingFactory{})
	return &settingClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type PublicDomainsGetter interface {
	PublicDomains(namespace string) PublicDomainInterface
}

func (c *Client) PublicDomains(namespace string) PublicDomainInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &PublicDomainResource, PublicDomainGroupVersionKind, publicDomainFactory{})
	return &publicDomainClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type FeaturesGetter interface {
	Features(namespace string) FeatureInterface
}

func (c *Client) Features(namespace string) FeatureInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &FeatureResource, FeatureGroupVersionKind, featureFactory{})
	return &featureClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

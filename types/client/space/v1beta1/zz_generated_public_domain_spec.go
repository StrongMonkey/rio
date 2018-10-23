package client

const (
	PublicDomainSpecType              = "publicDomainSpec"
	PublicDomainSpecFieldDomainName   = "domainName"
	PublicDomainSpecFieldRouteSetName = "routeSetName"
	PublicDomainSpecFieldServiceName  = "serviceName"
	PublicDomainSpecFieldSpaceName    = "spaceName"
	PublicDomainSpecFieldStackName    = "stackName"
)

type PublicDomainSpec struct {
	DomainName   string `json:"domainName,omitempty" yaml:"domainName,omitempty"`
	RouteSetName string `json:"routeSetName,omitempty" yaml:"routeSetName,omitempty"`
	ServiceName  string `json:"serviceName,omitempty" yaml:"serviceName,omitempty"`
	SpaceName    string `json:"spaceName,omitempty" yaml:"spaceName,omitempty"`
	StackName    string `json:"stackName,omitempty" yaml:"stackName,omitempty"`
}
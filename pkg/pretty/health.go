package pretty

import (
	"time"

	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/mapper"
	pm "github.com/rancher/rio/pkg/pretty/mapper"
	riov1 "github.com/rancher/rio/types/apis/rio.cattle.io/v1"
)

func health(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, riov1.HealthConfig{},
			pm.Shlex{Field: "test"},
			mapper.Move{From: "intervalSeconds", To: "interval"},
			pm.Duration{Field: "interval", Unit: time.Second},
			pm.AliasField{Field: "interval", Names: []string{"period", "periodSeconds"}},
			mapper.Move{From: "timeoutSeconds", To: "timeout"},
			pm.Duration{Field: "timeout", Unit: time.Second},
			mapper.Move{From: "initialDelaySeconds", To: "initialDelay"},
			pm.Duration{Field: "initialDelay", Unit: time.Second},
			pm.AliasField{Field: "initialDelay", Names: []string{"startPeriod"}},
			pm.AliasField{Field: "healthyThreshold", Names: []string{"retries", "successThreshold"}},
			pm.AliasField{Field: "unhealthyThreshold", Names: []string{"failureThreshold"}},
		)
}

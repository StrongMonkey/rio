package integration

import (
	"os"
	"testing"

	"github.com/rancher/rio/tests/testutil"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMain(m *testing.M) {
	testutil.IntegrationPreCheck()
	os.Exit(m.Run())
}

func TestSuite(t *testing.T) {
	suite := spec.New("integration suite", spec.Report(report.Terminal{}))
	specs := map[string]func(t *testing.T, when spec.G, it spec.S){
		"run":             runTests,      // base done
		"scale":           scaleTests,    // base done
		"weight":          weightTests,   // base done
		"endpoint":        endpointTests, // base done
		"domain":          domainTests,
		"route":           routeTests,
		"export":          exportTests, // base done
		"config":          configTests,
		"externalService": externalServiceTests,
		//"riofile":         riofileTests,
		//"rbac":            rbacTests,
	}
	for desc, fnc := range specs {
		suite(desc, fnc)
	}
	suite.Run(t)
}

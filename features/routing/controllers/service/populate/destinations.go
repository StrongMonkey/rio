package populate

import (
	"github.com/knative/pkg/apis/istio/v1alpha3"
	"github.com/rancher/norman/pkg/objectset"
	"github.com/rancher/rio/features/stack/controllers/service/populate/servicelabels"
	"github.com/rancher/rio/pkg/serviceset"
	v1alpha3client "github.com/rancher/rio/types/apis/networking.istio.io/v1alpha3"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DestinationRulesAndVirtualServices(stack *v1.Stack, services []*v1.Service, service *v1.Service, os *objectset.ObjectSet) error {
	if err := destinationRules(stack, services, service, os); err != nil {
		return err
	}
	return virtualServices(stack, services, service, os)
}

func destinationRules(stack *v1.Stack, services []*v1.Service, service *v1.Service, os *objectset.ObjectSet) error {
	serviceSets, err := serviceset.CollectionServices(services)
	if err != nil {
		return err
	}

	serviceSet, ok := serviceSets[service.Name]
	if !ok {
		return nil
	}

	dr := destinationRuleForService(stack, service.Name, serviceSet)
	os.Add(dr)

	for _, dr := range destinationRuleForRevisionService(stack, serviceSet) {
		os.Add(dr)
	}

	return nil
}

func destinationRuleForService(stack *v1.Stack, name string, service *serviceset.ServiceSet) *v1alpha3client.DestinationRule {
	drSpec := v1alpha3.DestinationRuleSpec{
		Host: name,
	}

	drSpec.Subsets = append(drSpec.Subsets, newSubSet(service.Service))

	for _, rev := range service.Revisions {
		drSpec.Subsets = append(drSpec.Subsets, newSubSet(rev))
	}

	dr := newDestinationRule(stack, service.Service)
	dr.Spec = drSpec

	return dr
}

func destinationRuleForRevisionService(stack *v1.Stack, service *serviceset.ServiceSet) []*v1alpha3client.DestinationRule {
	var result []*v1alpha3client.DestinationRule

	for _, rev := range service.Revisions {
		drObject := newDestinationRule(stack, rev)
		drObject.Spec = v1alpha3.DestinationRuleSpec{
			Host: rev.Name,
			Subsets: []v1alpha3.Subset{
				newSubSet(rev),
			},
		}
		result = append(result, drObject)
	}

	return result
}

func newSubSet(service *v1.Service) v1alpha3.Subset {
	return v1alpha3.Subset{
		Name: service.Spec.Revision.Version,
		Labels: map[string]string{
			"rio.cattle.io/version": service.Spec.Revision.Version,
		},
	}
}

func newDestinationRule(stack *v1.Stack, service *v1.Service) *v1alpha3client.DestinationRule {
	return v1alpha3client.NewDestinationRule(service.Namespace, service.Name, v1alpha3client.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Labels: servicelabels.RioOnlyServiceLabels(stack, service),
		},
	})
}

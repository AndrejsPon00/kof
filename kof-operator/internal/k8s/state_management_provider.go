package k8s

import (
	"context"

	kcmv1beta1 "github.com/K0rdent/kcm/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetStateManagementsProviders(ctx context.Context, client client.Client) (*kcmv1beta1.StateManagementProviderList, error) {
	smpList := &kcmv1beta1.StateManagementProviderList{
		Items: make([]kcmv1beta1.StateManagementProvider, 0),
	}
	err := client.List(ctx, smpList)
	return smpList, err
}

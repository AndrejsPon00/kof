package k8s

import (
	"context"

	kcmv1beta1 "github.com/K0rdent/kcm/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetMultiClusterServices(ctx context.Context, client client.Client) (*kcmv1beta1.MultiClusterServiceList, error) {
	mcsList := &kcmv1beta1.MultiClusterServiceList{
		Items: make([]kcmv1beta1.MultiClusterService, 0),
	}
	err := client.List(ctx, mcsList)
	return mcsList, err
}

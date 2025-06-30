package handlers

import (
	"context"
	"errors"
	"net/http"

	kcmv1beta1 "github.com/K0rdent/kcm/api/v1beta1"
	"github.com/go-logr/logr"
	"github.com/k0rdent/kof/kof-operator/internal/k8s"
	"github.com/k0rdent/kof/kof-operator/internal/server"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ObservabilityHandler struct {
	kubeClient *k8s.KubeClient
	logger     *logr.Logger
}

type Response struct {
	ClusterDeployments       []kcmv1beta1.ClusterDeployment       `json:"cluster_deployments"`
	MultiClusterServices     []kcmv1beta1.MultiClusterService     `json:"multi_cluster_services"`
	StateManagementProviders []kcmv1beta1.StateManagementProvider `json:"state_management_provider"`
}

type ObjectsCollection struct {
	ClusterDeployments       *kcmv1beta1.ClusterDeploymentList
	MultiClusterServices     *kcmv1beta1.MultiClusterServiceList
	StateManagementProviders *kcmv1beta1.StateManagementProviderList
}

type ObjectResult[T any] struct {
	Result T
	Err    error
}

func newObservabilityHandler(res *server.Response) (*ObservabilityHandler, error) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	return &ObservabilityHandler{
		kubeClient: kubeClient,
		logger:     res.Logger,
	}, nil
}

func HandleObservability(res *server.Response, req *http.Request) {
	ctx := req.Context()

	h, err := newObservabilityHandler(res)
	if err != nil {
		res.Logger.Error(err, "Failed to create ... handler")
		res.Fail(server.BasicInternalErrorMessage, http.StatusInternalServerError)
		return
	}

	objectsCollection, err := h.collectObjects(ctx)
	if err != nil {
		res.Logger.Error(err, "Failed to collect objects")
		res.Fail(server.BasicInternalErrorMessage, http.StatusInternalServerError)
		return
	}

	response := createResponse(objectsCollection)
	res.Send(response, http.StatusOK)
}

func (h *ObservabilityHandler) collectObjects(ctx context.Context) (*ObjectsCollection, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cdChan := collectObject(ctx, h.kubeClient.Client, k8s.GetClusterDeployments, cancel)
	mcsChan := collectObject(ctx, h.kubeClient.Client, k8s.GetMultiClusterServices, cancel)
	smpChan := collectObject(ctx, h.kubeClient.Client, k8s.GetStateManagementsProviders, cancel)

	cdRes := <-cdChan
	if cdRes.Err != nil {
		return nil, cdRes.Err
	}

	mcsRes := <-mcsChan
	if mcsRes.Err != nil {
		return nil, mcsRes.Err
	}

	smpRes := <-smpChan
	if smpRes.Err != nil {
		return nil, smpRes.Err
	}

	return &ObjectsCollection{
		ClusterDeployments:       cdRes.Result,
		MultiClusterServices:     mcsRes.Result,
		StateManagementProviders: smpRes.Result,
	}, nil
}

func collectObject[T any](
	ctx context.Context,
	client client.Client,
	fetchFunc func(context.Context, client.Client) (T, error),
	cancel context.CancelFunc,
) <-chan ObjectResult[T] {
	resultChan := make(chan ObjectResult[T], 1)

	go func() {
		defer close(resultChan)
		result, err := fetchFunc(ctx, client)
		switch {
		case errors.Is(err, context.Canceled):
			err = nil
		case err != nil:
			cancel()
		}
		resultChan <- ObjectResult[T]{Result: result, Err: err}
	}()

	return resultChan
}

func createResponse(objects *ObjectsCollection) *Response {
	response := &Response{}

	if objects.ClusterDeployments != nil && len(objects.ClusterDeployments.Items) > 0 {
		response.ClusterDeployments = objects.ClusterDeployments.Items
	}

	if objects.MultiClusterServices != nil && len(objects.MultiClusterServices.Items) > 0 {
		response.MultiClusterServices = objects.MultiClusterServices.Items
	}

	if objects.StateManagementProviders != nil && len(objects.StateManagementProviders.Items) > 0 {
		response.StateManagementProviders = objects.StateManagementProviders.Items
	}

	return response
}

package client

import (
	"context"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"time"
)

type Argo struct {
	client      workflow.WorkflowServiceClient
	pollingTime time.Duration
}

type Opts struct {
	pollingTime time.Duration
}

type CreateRequest struct {
	Namespace string             `json:"namespace" yaml:"namespace"`
	Workflow  *v1alpha1.Workflow `json:"workflow" yaml:"workflow"`
}

type GetRequest struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	Name      string `json:"name" yaml:"name"`
}

func New(client workflow.WorkflowServiceClient, opts Opts) *Argo {
	return &Argo{
		client:      client,
		pollingTime: opts.pollingTime,
	}
}

func NewFromArgoServer(url string, opts Opts) (*Argo, error) {
	_, client, err := apiclient.NewClientFromOpts(apiclient.Opts{
		ArgoServerOpts: apiclient.ArgoServerOpts{
			URL: url,
		},
	})

	if err != nil {
		return nil, err
	}

	return New(client.NewWorkflowServiceClient(), opts), nil
}

func (a *Argo) Client() workflow.WorkflowServiceClient {
	return a.client
}

func (a *Argo) CreateWorkflow(ctx context.Context, req CreateRequest) (*v1alpha1.Workflow, error) {
	return a.client.CreateWorkflow(ctx, &workflow.WorkflowCreateRequest{
		Namespace: req.Namespace,
		Workflow:  req.Workflow,
	})
}

func (a *Argo) GetWorkflow(ctx context.Context, req GetRequest) (*v1alpha1.Workflow, error) {
	return a.client.GetWorkflow(ctx, &workflow.WorkflowGetRequest{
		Name:      req.Name,
		Namespace: req.Namespace,
	})
}

func (a *Argo) WaitWorkflow(ctx context.Context, req GetRequest) (*v1alpha1.Workflow, error) {
	tk := time.NewTimer(a.pollingTime)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-tk.C:
			wf, err := a.GetWorkflow(ctx, req)
			if err != nil {
				return nil, err
			}

			if wf.Status.Phase.Completed() {
				tk.Stop()
				return wf, nil
			}
		}
	}
}

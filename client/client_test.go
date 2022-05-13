package client

import (
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow/mocks"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestArgo_CreateWorkflow(t *testing.T) {
	expectedWf := v1alpha1.Workflow{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-workflow-00",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowStatus{
			Phase: v1alpha1.WorkflowPending,
		},
	}

	client := &mocks.WorkflowServiceClient{}
	client.On("CreateWorkflow", mock.Anything, mock.Anything).Return(&expectedWf, nil)

	argo := New(client, Opts{PollingTime: 100 * time.Millisecond})

	wf, err := argo.CreateWorkflow(context.TODO(), CreateRequest{
		Namespace: "default",
		Workflow:  &expectedWf,
	})

	require.NoError(t, err)
	require.NotNil(t, wf)

	require.Equal(t, expectedWf.Name, wf.Name)
	require.Equal(t, expectedWf.Namespace, wf.Namespace)
	require.Equal(t, expectedWf.Status.Phase, wf.Status.Phase)
}

func TestArgo_GetWorkflow(t *testing.T) {
	expectedWf := v1alpha1.Workflow{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-workflow-00",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowStatus{
			Phase: v1alpha1.WorkflowPending,
		},
	}

	client := &mocks.WorkflowServiceClient{}

	client.On("GetWorkflow", mock.Anything, mock.Anything).Return(&expectedWf, nil)

	argo := New(client, Opts{PollingTime: 100 * time.Millisecond})

	wf, err := argo.GetWorkflow(context.TODO(), GetRequest{
		Namespace: "default",
		Name:      "test-job-00",
	})

	require.NoError(t, err)
	require.NotNil(t, wf)
	require.Equal(t, expectedWf.Name, wf.Name)
	require.Equal(t, expectedWf.Namespace, wf.Namespace)
	require.Equal(t, expectedWf.Status.Phase, wf.Status.Phase)
}

func TestArgo_WaitWorkflow(t *testing.T) {
	expectedWf := v1alpha1.Workflow{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-workflow-00",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowStatus{
			Phase: v1alpha1.WorkflowSucceeded,
		},
	}

	client := &mocks.WorkflowServiceClient{}
	client.On("GetWorkflow", mock.Anything, mock.Anything).Return(&expectedWf, nil)

	argo := New(client, Opts{
		PollingTime: 100 * time.Millisecond,
	})

	wf, err := argo.WaitWorkflow(context.TODO(), GetRequest{
		Namespace: "default",
		Name:      "test-job-00",
	})

	require.NoError(t, err)
	require.NotNil(t, wf)

	require.Equal(t, expectedWf.Name, wf.Name)
	require.Equal(t, expectedWf.Namespace, wf.Namespace)
	require.Equal(t, expectedWf.Status.Phase, wf.Status.Phase)
}

func TestArgo_WaitWorkflow_Pending(t *testing.T) {
	expectedWf := v1alpha1.Workflow{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-workflow-00",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowStatus{
			Phase: v1alpha1.WorkflowPending,
		},
	}

	client := &mocks.WorkflowServiceClient{}
	client.On("GetWorkflow", mock.Anything, mock.Anything).Return(&expectedWf, nil)

	argo := New(client, Opts{
		PollingTime: 100 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	_, err := argo.WaitWorkflow(ctx, GetRequest{
		Namespace: "default",
		Name:      "test-job-00",
	})

	require.Error(t, err)
}

func TestArgo_WaitWorkflow_SwitchPendingToRunning(t *testing.T) {
	expectedWf := v1alpha1.Workflow{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-workflow-00",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowStatus{
			Phase: v1alpha1.WorkflowPending,
		},
	}

	completedWf := v1alpha1.Workflow{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-workflow-00",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowStatus{
			Phase: v1alpha1.WorkflowSucceeded,
		},
	}

	client := &mocks.WorkflowServiceClient{}

	// after first call switch the jobs status phase in order to terminate the polling
	client.On("GetWorkflow", mock.Anything, mock.Anything).
		Return(&expectedWf, nil).
		Return(&completedWf, nil)

	argo := New(client, Opts{
		PollingTime: 100 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	_, err := argo.WaitWorkflow(ctx, GetRequest{
		Namespace: "default",
		Name:      "test-job-00",
	})

	require.NoError(t, err)
}

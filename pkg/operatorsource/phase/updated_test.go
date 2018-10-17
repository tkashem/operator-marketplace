package phase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	gomock "github.com/golang/mock/gomock"
	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	mocks "github.com/operator-framework/operator-marketplace/pkg/mocks/operatorsource_mocks"
	"github.com/operator-framework/operator-marketplace/pkg/operatorsource/phase"
)

// Use Case: Admin has changed the Spec to point to different endpoint or namespace.
// Expected Result: Current operator source should be purged and the next phase
// should be set to "Validating" so that reconciliation is triggered.
func TestReconcile_SpecHasChanged_ReconciliationTriggered(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.TODO()

	opsrcIn := helperNewOperatorSource("marketplace", "foo", v1alpha1.OperatorSourcePhaseSucceeded)
	opsrcWant := opsrcIn.DeepCopy()
	opsrcWant.Status = v1alpha1.OperatorSourceStatus{}

	nextPhaseWant := &phase.NextPhase{
		Phase:   v1alpha1.OperatorSourcePhaseValidating,
		Message: v1alpha1.GetOperatorSourcePhaseMessage(v1alpha1.OperatorSourcePhaseValidating),
	}

	datastore := mocks.NewDatastoreWriter(controller)
	kubeclient := mocks.NewKubeClient(controller)
	reconciler := phase.NewUpdatedEventReconciler(helperGetContextLogger(), datastore, kubeclient)

	// We expect the operator source to be removed from the datastore.
	csc := helperNewCatalogSourceConfig(opsrcIn.Namespace, getExpectedCatalogSourceConfigName(opsrcIn.Name))
	datastore.EXPECT().Remove(opsrcIn.GetUID())

	// We expect the associated CatalogConfigSource object to be deleted.
	kubeclient.EXPECT().Delete(csc)

	opsrcGot, nextPhaseGot, errGot := reconciler.Reconcile(ctx, opsrcIn)

	assert.NoError(t, errGot)
	assert.Equal(t, opsrcWant, opsrcGot)
	assert.Equal(t, nextPhaseWant, nextPhaseGot)
}

// Use Case: The associated CatalogSourceConfig object is not found while purging.
// Expected Result: NotFound error is ignored and the next phase should be set
// to "Validating" so that reconciliation is triggered.
func TestReconcile_CatalogSourceConfigNotFoubd_ErrorExpected(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.TODO()

	opsrcIn := helperNewOperatorSource("marketplace", "foo", v1alpha1.OperatorSourcePhaseSucceeded)
	opsrcWant := opsrcIn.DeepCopy()
	opsrcWant.Status = v1alpha1.OperatorSourceStatus{}

	nextPhaseWant := &phase.NextPhase{
		Phase:   v1alpha1.OperatorSourcePhaseValidating,
		Message: v1alpha1.GetOperatorSourcePhaseMessage(v1alpha1.OperatorSourcePhaseValidating),
	}

	datastore := mocks.NewDatastoreWriter(controller)
	kubeclient := mocks.NewKubeClient(controller)
	reconciler := phase.NewUpdatedEventReconciler(helperGetContextLogger(), datastore, kubeclient)

	// We expect the operator source to be removed from the datastore.
	csc := helperNewCatalogSourceConfig(opsrcIn.Namespace, getExpectedCatalogSourceConfigName(opsrcIn.Name))
	datastore.EXPECT().Remove(opsrcIn.GetUID())

	// We expect kube client to throw a NotFound error.
	notFoundErr := k8s_errors.NewNotFound(schema.GroupResource{}, "CatalogSourceConfig not found")
	kubeclient.EXPECT().Delete(csc).Return(notFoundErr)

	opsrcGot, nextPhaseGot, errGot := reconciler.Reconcile(ctx, opsrcIn)

	assert.Error(t, errGot)
	assert.Equal(t, opsrcGot, opsrcWant)
	assert.Equal(t, nextPhaseWant, nextPhaseGot)
}

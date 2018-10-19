package phase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	gomock "github.com/golang/mock/gomock"
	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/appregistry"
	mocks "github.com/operator-framework/operator-marketplace/pkg/mocks/operatorsource_mocks"
	"github.com/operator-framework/operator-marketplace/pkg/operatorsource/phase"
)

// Use Case: Successfully validated and scheduled for download.
// Expected Result: Manifest(s) downloaded and stored successfully and the next
// phase set to "Configuring".
func TestReconcile_ScheduledForDownload_Success(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	nextPhaseWant := &phase.NextPhase{
		Phase:   v1alpha1.OperatorSourcePhaseConfiguring,
		Message: v1alpha1.GetOperatorSourcePhaseMessage(v1alpha1.OperatorSourcePhaseConfiguring),
	}

	datastore := mocks.NewDatastoreWriter(controller)
	factory := mocks.NewAppRegistryClientFactory(controller)

	reconciler := phase.NewDownloadingReconciler(helperGetContextLogger(), factory, datastore)

	ctx := context.TODO()
	opsrcIn := helperNewOperatorSource("marketplace", "foo", v1alpha1.OperatorSourcePhaseDownloading)

	registryClient := mocks.NewAppRegistryClient(controller)
	factory.EXPECT().New(opsrcIn.Spec.Type, opsrcIn.Spec.Endpoint).Return(registryClient, nil).Times(1)

	// We expect the remote registry to return a non-empty list of manifest(s).
	manifestExpected := []*appregistry.OperatorMetadata{
		&appregistry.OperatorMetadata{
			Namespace:  "redhat",
			Repository: "myapp",
			Release:    "1.0.0",
			Digest:     "abcdefgh",
		},
	}
	registryClient.EXPECT().RetrieveAll(opsrcIn.Spec.RegistryNamespace).Return(manifestExpected, nil).Times(1)

	// We expect the datastore to save downloaded manifest(s) returned by the registry.
	datastore.EXPECT().Write(opsrcIn, manifestExpected).Return(nil)

	opsrcGot, nextPhaseGot, errGot := reconciler.Reconcile(ctx, opsrcIn)

	assert.NoError(t, errGot)
	assert.Equal(t, opsrcIn, opsrcGot)
	assert.Equal(t, nextPhaseWant, nextPhaseGot)
}

// Use Case: Registry returns an empty list of manifest(s).
// Expected Result: Next phase is set to "Failed".
func TestReconcile_OperatorSourceReturnsEmptyManifestList_ErrorExpected(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	datastore := mocks.NewDatastoreWriter(controller)
	factory := mocks.NewAppRegistryClientFactory(controller)

	reconciler := phase.NewDownloadingReconciler(helperGetContextLogger(), factory, datastore)

	ctx := context.TODO()
	opsrcIn := helperNewOperatorSource("marketplace", "foo", v1alpha1.OperatorSourcePhaseDownloading)

	registryClient := mocks.NewAppRegistryClient(controller)
	factory.EXPECT().New(opsrcIn.Spec.Type, opsrcIn.Spec.Endpoint).Return(registryClient, nil).Times(1)

	// We expect the registry to return an empty manifest list.
	manifests := []*appregistry.OperatorMetadata{}
	registryClient.EXPECT().RetrieveAll(opsrcIn.Spec.RegistryNamespace).Return(manifests, nil).Times(1)

	// Even through the registry returned an empty list, we expect the datastore
	// to save the operator source Spec along with the empty list.
	datastore.EXPECT().Write(opsrcIn, manifests).Return(nil)

	opsrcGot, nextPhaseGot, errGot := reconciler.Reconcile(ctx, opsrcIn)
	assert.Error(t, errGot)

	nextPhaseWant := &phase.NextPhase{
		Phase:   v1alpha1.OperatorSourcePhaseFailed,
		Message: errGot.Error(),
	}

	assert.Equal(t, opsrcIn, opsrcGot)
	assert.Equal(t, nextPhaseWant, nextPhaseGot)
}

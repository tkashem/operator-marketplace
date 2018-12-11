package operatorsource

import (
	"context"
	"errors"

	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/appregistry"
	"github.com/operator-framework/operator-marketplace/pkg/datastore"
	"github.com/operator-framework/operator-marketplace/pkg/phase"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type UpdateChecker interface {
	Check(types.UID) (bool, error)
	Trigger(types.NamespacedName) (deleted bool, updateErr error)
}

type checker struct {
	factory      appregistry.ClientFactory
	datastore    datastore.Writer
	client       client.Client
	transitioner phase.Transitioner
}

func (c *checker) Check(opsrcUID types.UID) (bool, error) {
	spec, exists := c.datastore.GetOperatorSource(opsrcUID)
	if !exists {
		return false, errors.New("The given OperatorSource object does not exist in datastore")
	}

	registry, err := c.factory.New(spec.Type, spec.Endpoint)
	if err != nil {
		return false, err
	}

	metadata, err := registry.ListPackages(spec.RegistryNamespace)
	if err != nil {
		return false, err
	}

	updated, err := c.datastore.HasUpdates(opsrcUID, metadata)
	if err != nil {
		return false, err
	}

	return updated, nil
}

func (c *checker) Trigger(namespacedName types.NamespacedName) (deleted bool, updateErr error) {
	instance := &v1alpha1.OperatorSource{}

	// Get the current state of the given object before we make any decision.
	if err := c.client.Get(context.TODO(), namespacedName, instance); err != nil {
		// Not found, the given OperatorSource object could have been deleted.
		// Treat it as no error and indicate that the object has been deleted.
		if k8s_errors.IsNotFound(err) {
			deleted = true
			return
		}

		// Otherwise it is an error
		updateErr = err
		return
	}

	// Needed because sdk does not get the gvk
	instance.EnsureGVK()

	// We want to purge the OperatorSource object so that the cache can rebuild.
	nextPhase := &v1alpha1.Phase{
		Name:    phase.OperatorSourcePurging,
		Message: "Remote registry has been updated",
	}
	if !c.transitioner.TransitionInto(&instance.Status.CurrentPhase, nextPhase) {
		// No need to update since the object is already in purging phase.
		return
	}

	if err := c.client.Update(context.TODO(), instance); err != nil {
		updateErr = err
		return
	}

	return
}

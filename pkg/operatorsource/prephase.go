package operatorsource

import (
	"context"

	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/datastore"
	"github.com/operator-framework/operator-marketplace/pkg/phase"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewPrePhaseReconciler returns a pre-phase Reconciler.
func NewPrePhaseReconciler(logger *log.Entry, datastore datastore.Writer, client client.Client) Reconciler {
	return &prePhaseReconciler{
		logger:    logger,
		datastore: datastore,
		client:    client,
	}
}

// prePhaseReconciler implements Reconciler interface.
type prePhaseReconciler struct {
	logger    *log.Entry
	datastore datastore.Writer
	client    client.Client
}

// Reconcile determines if the following conditions are true and
// act appropriately.
//
// A. An admin changes the spec of a given OperatorSource object. This
//    warrants for a purge and reconciliation to start anew.
// B. When marketplace operator restarts it loses the in-memory cache and
//    so it needs to rebuild the cache for all existing OperatorSource(s).
//
// in represents the original OperatorSource object received from the sdk
// and before reconciliation has started.
//
// out represents the OperatorSource object after reconciliation has completed
// and could be different from the original. The OperatorSource object received
// (in) should be deep copied into (out) before changes are made.
//
// nextPhase represents the next desired phase for the given OperatorSource
// object. If nil is returned, it implies that no phase transition is expected.
func (r *prePhaseReconciler) Reconcile(ctx context.Context, in *v1alpha1.OperatorSource) (out *v1alpha1.OperatorSource, nextPhase *v1alpha1.Phase, err error) {
	out = in.DeepCopy()

	currentPhase := in.GetCurrentPhaseName()

	// If the object is in Initial phase, this implies that it's a new
	// OperatorSource CR or the user dropped the status field of an existing CR.
	// In either case, bail out and let the regular phased reconciler handle
	// the specified OperatorSource object.
	// If the OperatorSource object is in "Purging" phase, then return.
	if currentPhase == phase.Initial || currentPhase == phase.OperatorSourcePurging {
		return
	}

	// If we are here, the object is in a particular phase. Now is the time to
	// determine if Spec of the OperatorSource object has been updated.
	oldSpec, exists := r.datastore.GetOperatorSource(in.GetUID())
	if exists && oldSpec.IsEqual(&in.Spec) {
		// Underlying datastore is aware of the OperatorSource object and
		// Spec has not been modified.
		// Let the regular phased reconciler handle the OperatorSource object.
		return
	}

	// If we are here, it implies the following two scenarios:
	//
	// A. The underlying datastore is aware of the OperatorSource object and
	//    the Spec has been modified.
	//
	// B. The underlying datastore is not aware of this OperatorSource object
	//    and it has a valid state. This implies that the cache is out of sync.
	//
	// In either case, we want to purge the OperatorSource object and kick off
	// phased reconciliation from Purging phase.
	nextPhase = phase.GetNext(phase.OperatorSourcePurging)

	r.logger.Info("Out of sync, scheduling for reconciliation from 'Purging' phase")

	return
}

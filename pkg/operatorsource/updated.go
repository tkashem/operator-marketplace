package operatorsource

import (
	"context"

	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/datastore"
	"github.com/operator-framework/operator-marketplace/pkg/phase"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewPrePhaseReconciler returns a Reconciler that reconciles
// an OperatorSource object that has been updated.
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

// Reconcile reconciles an OperatorSource object whose Spec has been updated.
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
//
// On an update we purge the current OperatorSource object, drop the Status
// field and trigger reconciliation anew from "Validating" phase.
//
// If the purge fails the OperatorSource object is moved to "Failed" phase.
func (r *prePhaseReconciler) Reconcile(ctx context.Context, in *v1alpha1.OperatorSource) (out *v1alpha1.OperatorSource, nextPhase *v1alpha1.Phase, err error) {
	out = in.DeepCopy()

	currentPhase := in.Status.CurrentPhase.Name

	// If the object is in Initial state, this implies that
	// It's a new OperatorSource CR or the user dropped the status field of an
	// existing CR.
	// In either case we want to kick of phased reconciliation.
	if currentPhase == phase.Initial {
		return
	}

	oldSpec, exists := r.datastore.GetOperatorSource(in.GetUID())
	if exists && oldSpec.IsEqual(&in.Spec) {
		return
	}

	// If the underlying datastore is not aware of this operator source
	// and it has a valid state, then the cache is out of sync.
	// It's time to rebuild the cache.

	// If the Spec of the given OperatorSource object has changed from
	// the one in datastore then we treat it as an update event.
	r.datastore.RemoveOperatorSource(in.GetUID())

	// Drop existing Status field so that reconciliation can start anew.
	out.Status = v1alpha1.OperatorSourceStatus{}

	nextPhase = phase.GetNext(phase.OperatorSourcePurging)
	r.logger.Info("Spec has changed, scheduling for reconciliation from 'Purging' phase")

	return
}

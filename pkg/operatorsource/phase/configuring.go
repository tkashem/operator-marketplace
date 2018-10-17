package phase

import (
	"context"

	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/datastore"
	"github.com/operator-framework/operator-marketplace/pkg/kube"
	log "github.com/sirupsen/logrus"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
)

// NewConfiguringReconciler returns a Reconciler that reconciles
// an OperatorSource object in "Configuring" phase.
func NewConfiguringReconciler(logger *log.Entry, datastore datastore.Writer, kubeclient kube.Client) Reconciler {
	return &configuringReconciler{
		logger:     logger,
		datastore:  datastore,
		kubeclient: kubeclient,
		builder:    &CatalogSourceConfigBuilder{},
	}
}

// configuringReconciler is an implementation of Reconciler interface that
// reconciles an OperatorSource object in "Configuring" phase.
type configuringReconciler struct {
	logger     *log.Entry
	datastore  datastore.Writer
	kubeclient kube.Client
	builder    *CatalogSourceConfigBuilder
}

// Reconcile reconciles an OperatorSource object that is in "Configuring" phase.
// It ensures that a corresponding CatalogSourceConfig object exists.
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
// Upon success, it returns "Succeeded" as the next and final desired phase.
// On error, the function returns "Failed" as the next desied phase
// and Message is set to appropriate error message.
//
// If the corresponding CatalogSourceConfig object already exists
// then no further action is taken.
func (r *configuringReconciler) Reconcile(ctx context.Context, in *v1alpha1.OperatorSource) (out *v1alpha1.OperatorSource, nextPhase *NextPhase, err error) {
	if in.Status.Phase != v1alpha1.OperatorSourcePhaseConfiguring {
		err = ErrWrongReconcilerInvoked
		return
	}

	out = in

	cscName := getCatalogSourceConfigName(in.Name)
	cscRetrievedInto := r.builder.WithMeta(in.Namespace, cscName).CatalogSourceConfig()

	err = r.kubeclient.Get(cscRetrievedInto)

	if err == nil {
		r.logger.Infof("No action taken, CatalogSourceConfig [name=%s] already exists", cscName)
		nextPhase = getNextPhase(v1alpha1.OperatorSourcePhaseSucceeded)
		return
	}

	if !k8s_errors.IsNotFound(err) {
		nextPhase = getNextPhaseWithMessage(v1alpha1.OperatorSourcePhaseFailed, err.Error())
		return
	}

	manifests := r.datastore.GetPackageIDs(in.GetUID())

	csc := r.builder.WithMeta(in.Namespace, cscName).
		WithSpec(in.Namespace, manifests).
		WithOwner(in).
		CatalogSourceConfig()

	err = r.kubeclient.Create(csc)
	if err != nil {
		nextPhase = getNextPhaseWithMessage(v1alpha1.OperatorSourcePhaseFailed, err.Error())
		return
	}

	nextPhase = getNextPhase(v1alpha1.OperatorSourcePhaseSucceeded)
	r.logger.Info("The object has been successfully reconciled")

	return
}

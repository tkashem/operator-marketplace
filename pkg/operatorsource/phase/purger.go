package phase

import (
	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/datastore"
	"github.com/operator-framework/operator-marketplace/pkg/kube"

	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
)

// NewPurger returns a new instance of Purger.
func NewPurger(datastore datastore.Writer, kubeclient kube.Client) Purger {
	return &purger{
		datastore:  datastore,
		kubeclient: kubeclient,
	}
}

// Purger is an interface that wraps the Purge method.
//
// Purge purges all artifacts associated with an OperatorSource object. Operator
// manifest(s) associated with the OperatorSource object are removed from the
// underlying datastore and the corresponding CatalogSourceConfig object is
// deleted.
type Purger interface {
	Purge(in *v1alpha1.OperatorSource) error
}

// purger implements the Purger interface.
type purger struct {
	datastore  datastore.Writer
	kubeclient kube.Client
}

func (p *purger) Purge(in *v1alpha1.OperatorSource) error {
	p.datastore.Remove(in.GetUID())

	builder := &CatalogSourceConfigBuilder{}
	csc := builder.WithMeta(in.Namespace, getCatalogSourceConfigName(in.Name)).CatalogSourceConfig()

	err := p.kubeclient.Delete(csc)
	if err != nil && k8s_errors.IsNotFound(err) {
		return nil
	}

	return err
}

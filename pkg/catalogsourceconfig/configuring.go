package catalogsourceconfig

import (
	"context"
	"fmt"
	"strings"

	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/datastore"
	"github.com/operator-framework/operator-marketplace/pkg/phase"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

// Copied from https://github.com/operator-framework/operator-lifecycle-manager/blob/master/pkg/controller/registry/configmap_loader.go#L18
// TBD: Vendor in the folder once we require more than just these constants from
// the OLM registry code.
const (
	ConfigMapCRDName     = "customResourceDefinitions"
	ConfigMapCSVName     = "clusterServiceVersions"
	ConfigMapPackageName = "packages"
)

// NewConfiguringReconciler returns a Reconciler that reconciles a
// CatalogSourceConfig object in the "Configuring" phase.
func NewConfiguringReconciler(log *logrus.Entry, reader datastore.Reader) Reconciler {
	return &configuringReconciler{
		log:    log,
		reader: reader,
	}
}

// configuringReconciler is an implementation of Reconciler interface that
// reconciles a CatalogSourceConfig object in the "Configuring" phase.
type configuringReconciler struct {
	log    *logrus.Entry
	reader datastore.Reader
}

// Reconcile reconciles a CatalogSourceConfig object that is in the
// "Configuring" phase. It ensures that a corresponding CatalogSource object
// exists.
//
// Upon success, it returns "Succeeded" as the next and final desired phase.
// On error, the function returns "Failed" as the next desired phase
// and Message is set to the appropriate error message.
func (r *configuringReconciler) Reconcile(ctx context.Context, in *v1alpha1.CatalogSourceConfig) (out *v1alpha1.CatalogSourceConfig, nextPhase *v1alpha1.Phase, err error) {
	if in.Status.CurrentPhase.Name != phase.Configuring {
		err = phase.ErrWrongReconcilerInvoked
		return
	}

	out = in

	err = r.createCatalogSource(in)
	if err != nil {
		nextPhase = phase.GetNextWithMessage(phase.Failed, err.Error())
		return
	}

	nextPhase = phase.GetNext(phase.Succeeded)

	r.log.Info("The object has been successfully reconciled")
	return
}

// createCatalogData constructs the ConfigMap data by reading the manifest
// information of all packages from the datasource.
func (r *configuringReconciler) createCatalogData(csc *v1alpha1.CatalogSourceConfig) (map[string]string, error) {
	packageIDs := getPackageIDs(csc.Spec.Packages)
	data := make(map[string]string)
	if len(packageIDs) < 1 {
		return data, fmt.Errorf("No packages specified in CatalogSourceConfig %s/%s", csc.Namespace, csc.Name)
	}

	manifest, err := r.reader.Read(packageIDs)
	if err != nil {
		log.Errorf("Error \"%v\" getting manifest", err)
		return nil, err
	}

	// TBD: Do we create a CatalogSource per package?
	// TODO: Add more error checking
	data[ConfigMapCRDName] += manifest.CustomResourceDefinitions
	data[ConfigMapCSVName] += manifest.ClusterServiceVersions
	data[ConfigMapPackageName] += manifest.Packages

	return data, nil
}

// createCatalogSource creates a new CatalogSource CR and all the resources it
// requires.
func (r *configuringReconciler) createCatalogSource(csc *v1alpha1.CatalogSourceConfig) error {
	// Construct the operator artifact data to be placed in the ConfigMap data
	// section.
	data, err := r.createCatalogData(csc)
	if err != nil {
		return err
	}

	// Create the ConfigMap that will be used by the CatalogSource.
	catalogConfigMap := newConfigMap(csc, data)
	r.log.Infof("Creating %s ConfigMap", catalogConfigMap.Name)
	err = sdk.Create(catalogConfigMap)
	if err != nil && !errors.IsAlreadyExists(err) {
		r.log.Errorf("Failed to create ConfigMap : %v", err)
		return err
	}
	r.log.Infof("Created ConfigMap %s", catalogConfigMap.Name)

	catalogSource := newCatalogSource(csc, catalogConfigMap.Name)
	err = sdk.Create(catalogSource)
	if err != nil && !errors.IsAlreadyExists(err) {
		r.log.Errorf("Failed to create CatalogSource : %v", err)
		return err
	}
	r.log.Infof("Created CatalogSource %s", catalogSource.Name)

	return nil
}

// getPackageIDs returns a list of IDs from a comma separated string of IDs.
func getPackageIDs(csIDs string) []string {
	return strings.Split(csIDs, ",")
}

// newConfigMap returns a new ConfigMap object.
func newConfigMap(csc *v1alpha1.CatalogSourceConfig, data map[string]string) *corev1.ConfigMap {
	return new(ConfigMapBuilder).
		WithMeta(v1alpha1.ConfigMapPrefix+csc.Name, csc.Spec.TargetNamespace).
		WithOwner(csc).
		WithData(data).
		ConfigMap()
}

// newCatalogSource returns a CatalogSource object.
func newCatalogSource(csc *v1alpha1.CatalogSourceConfig, configMapName string) *olm.CatalogSource {
	return new(CatalogSourceBuilder).
		WithMeta(v1alpha1.CatalogSourcePrefix+csc.Name, csc.Spec.TargetNamespace).
		WithOwner(csc).
		// TBD: where do we get display name and publisher from?
		WithSpec("internal", configMapName, csc.Name, csc.Name).
		CatalogSource()
}

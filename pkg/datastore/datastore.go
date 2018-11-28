package datastore

import (
	"errors"
	"fmt"
	"strings"

	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	Cache *memoryDatastore
)

func init() {
	// Cache is the global instance of datastore used by
	// the Marketplace operator.
	Cache = New()
}

// New returns an instance of memoryDatastore.
func New() *memoryDatastore {
	return &memoryDatastore{
		manifests: map[types.UID]*operatorSourceRow{},
		parser:    &manifestYAMLParser{},
		walker:    &walker{},
		bundler:   &bundler{},
	}
}

// Reader is the interface that wraps the Read method.
//
// Read prepares an operator manifest for a given set of package(s)
// uniquely represeneted by the package ID(s) specified in packageIDs. It
// returns an instance of RawOperatorManifestData that has the required set of
// CRD(s), CSV(s) and package(s).
//
// The manifest returned can be used by the caller to create a ConfigMap object
// for a catalog source (CatalogSource) in OLM.
type Reader interface {
	Read(packageIDs []string) (marshaled *RawOperatorManifestData, err error)
}

// Writer is an interface that is used to manage the underlying datastore
// for operator manifest.
type Writer interface {
	// GetPackageIDs returns a comma separated list of operator ID(s). Each ID
	// returned can be used to retrieve the manifest associated with the
	// operator from underlying datastore.
	GetPackageIDs() string

	// Write saves the Spec associated with a given OperatorSource object and
	// the downloaded manifest(s) into datastore.
	//
	// opsrc represents the given OperatorSource object.
	// manifests is the list of manifest(s) associated with a given operator
	// source.
	Write(opsrc *v1alpha1.OperatorSource, rawManifests []*OperatorMetadata) error
}

// operatorSourceRow is what gets stored in datastore after an OperatorSource CR
// is reconciled.
//
// Every reconciled OperatorSource object has a corresponding operatorSourceRow
// in datastore. The Writer interface accepts a raw operator manifest and
// marshals it into this type before writing it to the underlying storage.
type operatorSourceRow struct {
	// We store the Spec associated with a given OperatorSource object. This is
	// so that we can determine whether Spec for an existing operator source
	// has been updated.
	//
	// We compare the Spec of the received OperatorSource object to the one
	// in datastore.
	Spec *v1alpha1.OperatorSourceSpec

	Operators map[string]*SingleOperatorManifest
}

// memoryDatastore is an in-memory implementation of operator manifest datastore.
// TODO: In future, it will be replaced by an indexable persistent datastore.
type memoryDatastore struct {
	manifests map[types.UID]*operatorSourceRow
	parser    ManifestYAMLParser
	walker    ManifestWalker
	bundler   Bundler
}

func (ds *memoryDatastore) Read(packageIDs []string) (*RawOperatorManifestData, error) {
	manifests, err := ds.validate(packageIDs)
	if err != nil {
		return nil, err
	}

	manifest, err := ds.bundler.Bundle(manifests)
	if err != nil {
		return nil, fmt.Errorf("error while bundling package(s) into  manifest - %s", err)
	}

	return ds.parser.Marshal(manifest)
}

func (ds *memoryDatastore) Write(opsrc *v1alpha1.OperatorSource, rawManifests []*OperatorMetadata) error {
	if opsrc == nil || rawManifests == nil {
		return errors.New("invalid argument")
	}

	operators := map[string]*SingleOperatorManifest{}
	for _, rawManifest := range rawManifests {
		data, err := ds.parser.Unmarshal(rawManifest.RawYAML)
		if err != nil {
			return err
		}

		decomposer := newDecomposer()
		if err := ds.walker.Walk(data, decomposer); err != nil {
			return err
		}

		packages := decomposer.Packages()
		for i, operatorPackage := range packages {
			operators[operatorPackage.GetPackageID()] = packages[i]
		}
	}

	manifest := &operatorSourceRow{
		Spec:      &opsrc.Spec,
		Operators: operators,
	}

	ds.manifests[opsrc.GetUID()] = manifest

	return nil
}

func (ds *memoryDatastore) GetPackageIDs() string {
	keys := ds.getAllPackages()
	return strings.Join(keys, ",")
}

func (ds *memoryDatastore) getAllPackages() []string {
	keys := make([]string, 0)
	for _, row := range ds.manifests {
		for packageID, _ := range row.Operators {
			keys = append(keys, packageID)
		}
	}

	return keys
}

func (ds *memoryDatastore) getAllPackagesToMap() map[string]*SingleOperatorManifest {
	packages := map[string]*SingleOperatorManifest{}
	for _, row := range ds.manifests {
		for packageID, manifest := range row.Operators {
			packages[packageID] = manifest
		}
	}

	return packages
}

// validate ensures that no package is mentioned more than once in the list.
// It also ensures that the package(s) specified in the list has a corresponding
// manifest in the underlying datastore.
func (ds *memoryDatastore) validate(packageIDs []string) ([]*SingleOperatorManifest, error) {
	packages := make([]*SingleOperatorManifest, 0)
	packageMap := map[string]*SingleOperatorManifest{}

	existing := ds.getAllPackagesToMap()

	for _, packageID := range packageIDs {
		operatorPackage, exists := existing[packageID]
		if !exists {
			return nil, fmt.Errorf("package [%s] not found", packageID)
		}

		if _, exists := packageMap[packageID]; exists {
			return nil, fmt.Errorf("package [%s] has been specified more than once", packageID)
		}

		packageMap[packageID] = operatorPackage
		packages = append(packages, operatorPackage)
	}

	return packages, nil
}

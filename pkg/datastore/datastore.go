package datastore

import (
	"errors"
	"strings"

	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/appregistry"
	"k8s.io/apimachinery/pkg/types"
)

var (
	ErrManifestNotFound = errors.New("manifest not found")
)

// New returns a new instance of datastore for Operator Manifest(s)
func New() *memoryDatastore {
	return &memoryDatastore{
		rows: map[types.UID]*row{},
	}
}

// Reader is the interface that wraps the Read method
//
// Read returns the associated operator manifest given a package ID
type Reader interface {
	Read(packageID string) (*appregistry.OperatorMetadata, error)
}

// Writer is an interface that is used to manage the underlying datastore
// for operator manifest.
type Writer interface {
	// GetPackageIDs returns a comma separated list of IDs of
	// all package(s) from the underlying datastore for a given operator source.
	//
	// opsrcUID is the unique identifier associated with a given operator source.
	//
	// If the operator source specified by opsrcUID is not found in datastore
	// then an empty string is returned.
	GetPackageIDs(opsrcUID types.UID) string

	// Write saves the Spec associated with a given OperatorSource object and
	// the downloaded manifest(s) into datastore.
	//
	// opsrc represents the given OperatorSource object.
	// manifests is the list of manifest(s) associated with a given operator
	// source.
	Write(opsrc *v1alpha1.OperatorSource, manifests []*appregistry.OperatorMetadata) error
}

// row encapsulates what we store for each operator source.
// An OperatorSource object has an one to one mapping to a row in datastore.
type row struct {
	// We store the Spec associated with a given OperatorSource object. This is
	// so that we can determine whether Spec for an existing operator source
	// has been updated.
	//
	// We compare the Spec of the received OperatorSource object to the one
	// in datastore.
	Spec *v1alpha1.OperatorSourceSpec

	// This is the list of manifest(s) associated with a given operator source.
	Manifests map[string]*appregistry.OperatorMetadata
}

// memoryDatastore is an in-memory implementation of operator manifest datastore.
// TODO: In future, it will be replaced by an indexable persistent datastore.
type memoryDatastore struct {
	rows map[types.UID]*row
}

func (ds *memoryDatastore) Read(packageID string) (*appregistry.OperatorMetadata, error) {
	for _, row := range ds.rows {
		manifest, exists := row.Manifests[packageID]
		if exists {
			return manifest, nil
		}
	}

	return nil, ErrManifestNotFound
}

func (ds *memoryDatastore) Write(opsrc *v1alpha1.OperatorSource, packages []*appregistry.OperatorMetadata) error {
	if opsrc == nil || packages == nil {
		return errors.New("invalid argument")
	}

	manifests := map[string]*appregistry.OperatorMetadata{}
	for _, pkg := range packages {
		manifests[pkg.ID()] = pkg
	}

	ds.rows[opsrc.GetUID()] = &row{
		Spec:      &opsrc.Spec,
		Manifests: manifests,
	}

	return nil
}

func (ds *memoryDatastore) GetPackageIDs(opsrcUID types.UID) string {
	row, exists := ds.rows[opsrcUID]
	if !exists {
		return ""
	}

	keys := make([]string, 0, len(row.Manifests))
	for key := range row.Manifests {
		keys = append(keys, key)
	}

	return strings.Join(keys, ",")
}

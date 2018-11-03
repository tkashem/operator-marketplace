package datastore

import (
	"errors"
	"strings"
)

var (
	ErrManifestNotFound = errors.New("manifest not found")
)

// New returns a new instance of datastore for Operator Manifest(s)
func New() *memoryDatastore {
	return &memoryDatastore{
		manifests:   map[string]*OperatorManifest{},
		unmarshaler: &blobUnmarshalerImpl{},
	}
}

// Reader is the interface that wraps the Read method
//
// Read returns an operator manifest for the given set packages represeneted
// by package ID(s) provided.
type Reader interface {
	Read(packageIDs []string) (*OperatorManifestData, error)
}

// Writer is an interface that is used to manage the underlying datastore
// for operator manifest.
type Writer interface {
	// GetPackageIDs returns a comma separated list of IDs of
	// all package(s) from underlying datastore.
	GetPackageIDs() string

	// Write stores the list of operator manifest(s) into datastore
	Write(packages []*OperatorMetadata) error
}

// memoryDatastore is an in-memory implementation of operator manifest datastore.
// TODO: In future, it will be replaced by an indexable persistent datastore.
type memoryDatastore struct {
	manifests   map[string]*OperatorManifest
	unmarshaler blobUnmarshaler
}

func (ds *memoryDatastore) Read(packageIDs []string) (*OperatorManifestData, error) {
	pkg := RawPackageData{}
	for _, packageID := range packageIDs {
		manifest, exists := ds.manifests[packageID]
		if !exists {
			return nil, ErrManifestNotFound
		}

		o, err := ds.unmarshaler.UnmarshalData(&manifest.Data)
		if err != nil {
			return nil, err
		}

		pkg.CustomResourceDefinitions = append(pkg.CustomResourceDefinitions, o.CustomResourceDefinitions...)
		pkg.ClusterServiceVersions = append(pkg.ClusterServiceVersions, o.ClusterServiceVersions...)
		pkg.Packages = append(pkg.Packages, o.Packages...)

	}

	return ds.unmarshaler.ToManifest(&pkg)
}

func (ds *memoryDatastore) Write(packages []*OperatorMetadata) error {
	for _, pkg := range packages {
		data, err := ds.unmarshaler.Unmarshal(pkg.RawYAML)
		if err != nil {
			return err
		}

		manifest := &OperatorManifest{
			RegistryMetadata: pkg.RegistryMetadata,
			Data:             data.Data,
		}

		ds.manifests[pkg.ID()] = manifest
	}

	return nil
}

func (ds *memoryDatastore) GetPackageIDs() string {
	keys := make([]string, 0, len(ds.manifests))
	for key := range ds.manifests {
		keys = append(keys, key)
	}

	return strings.Join(keys, ",")
}

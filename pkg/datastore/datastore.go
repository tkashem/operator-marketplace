package datastore

import (
	"fmt"
	"strings"
)

// New returns a new instance of datastore for Operator Manifest(s).
func New() *memoryDatastore {
	return &memoryDatastore{
		manifests: map[string]*OperatorManifest{},
		packages:  map[string]*ManifestPackage{},

		parser:   &manifestYAMLParser{},
		unpacker: &manifestUnpacker{},
	}
}

// Reader is the interface that wraps the Read method.
//
// Read prepares an operator manifest for a given set of package(s)
// uniquely represeneted by the package ID(s) specified in packageIDs. It
// returns an instance of OperatorManifestData that has the required set of
// CRD(s), CSV(s) and package(s).
//
// The manifest returned can be used by the caller to create a configMap object
// for a catalog source in OLM.
type Reader interface {
	Read(packageIDs []string) (marshaled *OperatorManifestData, err error)
}

// Writer is an interface that is used to manage the underlying datastore
// for operator manifest.
type Writer interface {
	// GetPackageIDs returns a comma separated list of IDs of
	// all package(s) from underlying datastore.
	GetPackageIDs() string

	// Write stores the list of operator manifest(s) into datastore.s
	Write(packages []*OperatorMetadata) error
}

// memoryDatastore is an in-memory implementation of operator manifest datastore.
// TODO: In future, it will be replaced by an indexable persistent datastore.
type memoryDatastore struct {
	manifests map[string]*OperatorManifest
	packages  map[string]*ManifestPackage

	parser   ManifestYAMLParser
	unpacker *manifestUnpacker
}

func (ds *memoryDatastore) Read(packageIDs []string) (*OperatorManifestData, error) {
	data := &Manifest{}
	for _, packageID := range packageIDs {
		operatorPackage, exists := ds.packages[packageID]
		if !exists {
			return nil, fmt.Errorf("package [%s] not found", packageID)
		}

		data.CustomResourceDefinitions = append(data.CustomResourceDefinitions, operatorPackage.CustomResourceDefinitions...)
		data.ClusterServiceVersions = append(data.ClusterServiceVersions, operatorPackage.ClusterServiceVersions...)
		data.Packages = append(data.Packages, operatorPackage.Package)
	}

	return ds.parser.Marshal(data)
}

func (ds *memoryDatastore) Write(packages []*OperatorMetadata) error {
	for _, pkg := range packages {
		data, err := ds.parser.Unmarshal(pkg.RawYAML)
		if err != nil {
			return err
		}

		packages, err := ds.unpacker.Unpack(data)
		if err != nil {
			return err
		}

		for _, operatorPackage := range packages {
			ds.packages[operatorPackage.Package.PackageName] = operatorPackage
		}

		manifest := &OperatorManifest{
			RegistryMetadata: pkg.RegistryMetadata,
			Data:             *data,
		}

		ds.manifests[pkg.ID()] = manifest
	}

	return nil
}

func (ds *memoryDatastore) GetPackageIDs() string {
	keys := make([]string, 0, len(ds.packages))
	for key := range ds.packages {
		keys = append(keys, key)
	}

	return strings.Join(keys, ",")
}

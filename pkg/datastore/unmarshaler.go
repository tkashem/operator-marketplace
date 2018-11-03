package datastore

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

type blobUnmarshaler interface {
	// Unmarshal unmarshals package blob into structured representations
	Unmarshal(in []byte) (*Manifest, error)

	UnmarshalData(data *OperatorManifestData) (*RawPackageData, error)

	ToManifest(*RawPackageData) (*OperatorManifestData, error)
}

type blobUnmarshalerImpl struct{}

func (*blobUnmarshalerImpl) Unmarshal(in []byte) (*Manifest, error) {
	data := &Manifest{}
	if err := yaml.Unmarshal(in, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (*blobUnmarshalerImpl) UnmarshalData(data *OperatorManifestData) (*RawPackageData, error) {
	crdJSON, err := yaml.YAMLToJSON([]byte(data.CustomResourceDefinitions))
	if err != nil {
		return nil, fmt.Errorf("error converting 'customResourceDefinitions' to JSON in package manifest : %s", err)
	}
	var crdList []KubeObject
	if err := json.Unmarshal(crdJSON, &crdList); err != nil {
		return nil, fmt.Errorf("error parsing 'customResourceDefinitions' in package manifest : %s", err)
	}

	csvJSON, err := yaml.YAMLToJSON([]byte(data.ClusterServiceVersions))
	if err != nil {
		return nil, fmt.Errorf("error converting 'clusterServiceVersions' to JSON in package manifest : %s", err)
	}
	var csvList []KubeObject
	if err := json.Unmarshal(csvJSON, &csvList); err != nil {
		return nil, fmt.Errorf("error parsing 'clusterServiceVersions' in package manifest : %s", err)
	}

	packageJSON, err := yaml.YAMLToJSON([]byte(data.Packages))
	if err != nil {
		return nil, fmt.Errorf("error converting 'packages' to JSON in package manifest : %s", err)
	}
	var packages []PackageManifest
	if err := json.Unmarshal(packageJSON, &packages); err != nil {
		return nil, fmt.Errorf("error parsing 'packages' in package manifest : %s", err)
	}

	return &RawPackageData{
		CustomResourceDefinitions: crdList,
		ClusterServiceVersions:    csvList,
		Packages:                  packages,
	}, nil
}

func (*blobUnmarshalerImpl) ToManifest(data *RawPackageData) (*OperatorManifestData, error) {
	crdRaw, err := yaml.Marshal(data.CustomResourceDefinitions)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'customResourceDefinitions' in package manifest : %s", err)
	}

	csvRaw, err := yaml.Marshal(data.ClusterServiceVersions)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'clusterServiceVersions' in package manifest : %s", err)
	}

	packageRaw, err := yaml.Marshal(data.Packages)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'packages' in package manifest : %s", err)
	}

	return &OperatorManifestData{
		CustomResourceDefinitions: string(crdRaw),
		ClusterServiceVersions:    string(csvRaw),
		Packages:                  string(packageRaw),
	}, nil
}

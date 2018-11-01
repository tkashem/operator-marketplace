package datastore

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Manifest encapsulates operator manifest data.
type Manifest struct {
	// Publisher represents the publisher of this package.
	Publisher string `yaml:"publisher"`

	// Data reflects the content of the package manifest.
	Data OperatorManifestData `yaml:"data"`
}

// OperatorManifestData encapsulates the list of CRD(s), CSV(s) and package(s)
// associated with an operator manifest.
type OperatorManifestData struct {
	// CustomResourceDefinitions is the set of custom resource definition(s)
	// associated with this package manifest.
	CustomResourceDefinitions string `yaml:"customResourceDefinitions"`

	// ClusterServiceVersions is the set of cluster service version(s)
	// associated with this package manifest.
	ClusterServiceVersions string `yaml:"clusterServiceVersions"`

	// Packages is the set of package(s) associated with this operator manifest.
	Packages string `yaml:"packages"`
}

// Manifest encapsulates operator manifest data
type OperatorManifest struct {
	// Metadata that uniquely identifies the given operator manifest in registry.
	RegistryMetadata RegistryMetadata

	// Data reflects the content of the package manifest.
	Data OperatorManifestData
}

// An operator manifest has three sections
// - customResourceDefinitions
// - clusterServiceVersions
// - packages
//
// This type is used to unmarshal an operator manifest into structured elements
// so that we can inspect a particular package, CustomResourceDefinition or
// ClusterServiceVersion.
//
type RawPackageData struct {
	// Set of custom resource definitions associated with this package manifest.
	CustomResourceDefinitions []KubeObject `json:"customResourceDefinitions"`

	// Set of cluster service version(s) associated with thus package manifest.
	ClusterServiceVersions []KubeObject `json:"clusterServiceVersions"`

	// Set of package(s) associated with this operator manifest.
	Packages []PackageManifest `json:"packages"`
}

// This type is used to unmarshal CustomResourceDefinition, ClusterServiceVersion
// from raw package manifest YAML.
//
// This allows us to achieve a loose coupling between marketplace and OLM.
// We don't want to put ourselves into a position where we have to rebuild our
// operator every time OLM decides to add or rename elements in either
// CustomResourceDefinition or ClusterServiceVersion.
type KubeObject struct {
	// Type metadata
	metav1.TypeMeta `json:",inline"`

	// Object metadata
	metav1.ObjectMeta `json:"metadata"`

	// We don't parse the spec field of CustomResourceDefinition or
	// ClusterServiceVersion since we are not interested in the content of spec
	// (at least for now)
	Spec json.RawMessage `json:"spec"`
}

// The following type has been copied as is from OLM
// https://github.com/operator-framework/operator-lifecycle-manager/blob/724b209ccfff33b6208cc5d05283800d6661d441/pkg/controller/registry/types.go#L78:6
//
// We use it to unmarshal 'packages' section of operator manifest package.
//
//
// PackageManifest holds information about a package, which is a reference to one (or more)
// channels under a single package.
type PackageManifest struct {
	// PackageName is the name of the overall package, ala `etcd`.
	PackageName string `json:"packageName"`

	// Channels are the declared channels for the package, ala `stable` or `alpha`.
	Channels []PackageChannel `json:"channels"`

	// DefaultChannelName is, if specified, the name of the default channel for the package. The
	// default channel will be installed if no other channel is explicitly given. If the package
	// has a single channel, then that channel is implicitly the default.
	DefaultChannelName string `json:"defaultChannel"`
}

// The following type has been directly copied as is from OLM
// https://github.com/operator-framework/operator-lifecycle-manager/blob/724b209ccfff33b6208cc5d05283800d6661d441/pkg/controller/registry/types.go#L105
//
// We use it to unmarshal 'channels' of operator manifest package.
//
// PackageChannel defines a single channel under a package, pointing to a
// version of that package.
type PackageChannel struct {
	// Name is the name of the channel, e.g. `alpha` or `stable`
	Name string `json:"name"`

	// CurrentCSVName defines a reference to the CSV holding the version of
	// this package currently for the channel.
	CurrentCSVName string `json:"currentCSV"`
}

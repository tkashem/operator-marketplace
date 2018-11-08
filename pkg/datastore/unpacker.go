package datastore

import (
	"errors"
	"fmt"

	olm_v1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

type ManifestPackage struct {
	Package                   *PackageManifest
	CustomResourceDefinitions []*v1beta1.CustomResourceDefinition
	ClusterServiceVersions    []*olm_v1alpha1.ClusterServiceVersion
}

type Manifest struct {
	Packages                  []*PackageManifest
	CustomResourceDefinitions []*v1beta1.CustomResourceDefinition
	ClusterServiceVersions    []*olm_v1alpha1.ClusterServiceVersion
}

type manifestUnpacker struct {
}

func (u *manifestUnpacker) Unpack(manifest *StructuredOperatorManifestData) ([]*ManifestPackage, error) {
	if manifest == nil {
		return nil, errors.New("manifest can not be <nil>")
	}

	packages := make([]*ManifestPackage, 0)

	crds := extractCustomResourceDefinitionsToMap(manifest)
	csvs := extractClusterServiceVersionsToMap(manifest)

	for i, p := range manifest.Packages {
		customResourceDefinitions, clusterServiceVersions, err := extract(&p, crds, csvs)
		if err != nil {
			return nil, fmt.Errorf("failed to extract CRD(s) for package [%s]", p.PackageName)
		}

		mp := &ManifestPackage{
			Package:                   &manifest.Packages[i],
			CustomResourceDefinitions: customResourceDefinitions,
			ClusterServiceVersions:    clusterServiceVersions,
		}

		packages = append(packages, mp)
	}

	return packages, nil
}

func extract(pm *PackageManifest, crds map[string]*v1beta1.CustomResourceDefinition, csvs map[string]*olm_v1alpha1.ClusterServiceVersion) ([]*v1beta1.CustomResourceDefinition, []*olm_v1alpha1.ClusterServiceVersion, error) {
	csvForThisPackage := make([]*olm_v1alpha1.ClusterServiceVersion, 0)
	crdForThisPackage := make([]*v1beta1.CustomResourceDefinition, 0)

	crdMapForPackage := map[string]*v1beta1.CustomResourceDefinition{}
	csvMapForPackage := map[string]*olm_v1alpha1.ClusterServiceVersion{}

	for _, channel := range pm.Channels {
		crdForCurrentChannel, csvForCurrentChannel, err := extractChannel(&channel, crds, csvs)
		if err != nil {
			return nil, nil, fmt.Errorf("error extracting package[%s] from manifest - %s", pm.PackageName, err)
		}

		crdUniqueForChannel := make([]*v1beta1.CustomResourceDefinition, 0)
		for _, crd := range crdForCurrentChannel {
			if _, ok := crdMapForPackage[crd.Name]; ok != true {
				crdMapForPackage[crd.Name] = crd
				crdUniqueForChannel = append(crdUniqueForChannel, crd)
			}
		}

		csvUniqueForChannel := make([]*olm_v1alpha1.ClusterServiceVersion, 0)
		for _, csv := range csvForCurrentChannel {
			if _, ok := csvMapForPackage[csv.Name]; ok != true {
				csvMapForPackage[csv.Name] = csv
				csvUniqueForChannel = append(csvUniqueForChannel, csv)
			}
		}

		csvForThisPackage = append(csvForThisPackage, csvUniqueForChannel...)
		crdForThisPackage = append(crdForThisPackage, crdUniqueForChannel...)
	}

	return crdForThisPackage, csvForThisPackage, nil
}

func extractChannel(channel *PackageChannel, crds map[string]*v1beta1.CustomResourceDefinition, csvs map[string]*olm_v1alpha1.ClusterServiceVersion) ([]*v1beta1.CustomResourceDefinition, []*olm_v1alpha1.ClusterServiceVersion, error) {
	csvForThisChannel := make([]*olm_v1alpha1.ClusterServiceVersion, 0)

	currentCSV, ok := csvs[channel.CurrentCSVName]
	if !ok {
		return nil, nil, fmt.Errorf("did not find CSV [%s] for channel [%s]", channel.CurrentCSVName, channel.Name)
	}

	olderCSVs, err := extractReplaces(currentCSV, csvs)
	if !ok {
		return nil, nil, fmt.Errorf("error finding CSV for channel [%s] - %s", channel.Name, err)
	}

	csvForThisChannel = append(csvForThisChannel, currentCSV)
	csvForThisChannel = append(csvForThisChannel, olderCSVs...)

	// Find all CRD(s) associated with this package
	crdForthisChannel, err := extractCRD(crds, csvForThisChannel...)
	if err != nil {
		return nil, nil, fmt.Errorf("error finding CRD for channel [%s] - %s", channel.Name, err)
	}

	return crdForthisChannel, csvForThisChannel, nil
}

func extractCRD(allCRDs map[string]*v1beta1.CustomResourceDefinition, csvs ...*olm_v1alpha1.ClusterServiceVersion) ([]*v1beta1.CustomResourceDefinition, error) {
	crds := make([]*v1beta1.CustomResourceDefinition, 0)

	for _, csv := range csvs {
		for _, owned := range csv.Spec.CustomResourceDefinitions.Owned {
			ownedCRD, ok := allCRDs[owned.Name]
			if !ok {
				return nil, fmt.Errorf("did not find a 'Owned' CRD[%s]", owned.Name)
			}

			crds = append(crds, ownedCRD)
		}

		for _, required := range csv.Spec.CustomResourceDefinitions.Required {
			requiredCRD, ok := allCRDs[required.Name]
			if !ok {
				return nil, fmt.Errorf("did not find a 'Required' CRD[%s]", required.Name)
			}

			crds = append(crds, requiredCRD)
		}
	}

	return crds, nil
}

func extractReplaces(currentCSV *olm_v1alpha1.ClusterServiceVersion, allCSVs map[string]*olm_v1alpha1.ClusterServiceVersion) ([]*olm_v1alpha1.ClusterServiceVersion, error) {
	olderCSVs := make([]*olm_v1alpha1.ClusterServiceVersion, 0)

	csv := currentCSV
	for {
		// If this CSV replaces an older version then we will include the older version as well.
		if csv.Spec.Replaces == "" {
			return olderCSVs, nil
		}

		replaceCSV, ok := allCSVs[csv.Spec.Replaces]
		if !ok {
			return nil, fmt.Errorf("did not find a 'Replaces' CSV [%s]", csv.Spec.Replaces)
		}

		olderCSVs = append(olderCSVs, replaceCSV)
		csv = replaceCSV
	}
}

func extractClusterServiceVersionsToMap(manifest *StructuredOperatorManifestData) map[string]*olm_v1alpha1.ClusterServiceVersion {
	csvs := map[string]*olm_v1alpha1.ClusterServiceVersion{}

	for i, csv := range manifest.ClusterServiceVersions {
		csvs[csv.Name] = &manifest.ClusterServiceVersions[i]
	}

	return csvs
}

func extractCustomResourceDefinitionsToMap(manifest *StructuredOperatorManifestData) map[string]*v1beta1.CustomResourceDefinition {
	crds := map[string]*v1beta1.CustomResourceDefinition{}

	for i, crd := range manifest.CustomResourceDefinitions {
		crds[crd.Name] = &manifest.CustomResourceDefinitions[i]
	}

	return crds
}

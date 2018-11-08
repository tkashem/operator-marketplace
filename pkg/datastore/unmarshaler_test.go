package datastore

import (
	"testing"

	olm_v1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// Do not use tabs for indentation as yaml forbids tabs http://yaml.org/faq.html
	rawCRDs = `
data:
  customResourceDefinitions: |-      
    - apiVersion: apiextensions.k8s.io/v1beta1
      kind: CustomResourceDefinition
      metadata:
        name: jbossapps-1.jboss.middleware.redhat.com
    - apiVersion: apiextensions.k8s.io/v1beta1
      kind: CustomResourceDefinition
      metadata:
        name: jbossapps-2.jboss.middleware.redhat.com
`

	crdWant = []*v1beta1.CustomResourceDefinition{
		&v1beta1.CustomResourceDefinition{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "jbossapps-1.jboss.middleware.redhat.com",
			},
		},
		&v1beta1.CustomResourceDefinition{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "jbossapps-2.jboss.middleware.redhat.com",
			},
		},
	}

	// Do not use tabs for indentation as yaml forbids tabs http://yaml.org/faq.html
	rawPackages = `
data:
  packages: |-
    - #! package-manifest: ./deploy/chart/catalog_resources/rh-operators/etcdoperator.v0.9.2.clusterserviceversion.yaml
      packageName: etcd
      channels:
        - name: alpha
          currentCSV: etcdoperator.v0.9.2
        - name: nightly
          currentCSV: etcdoperator.v0.9.2-nightly
      defaultChannel: alpha
`

	packagesWant = []*PackageManifest{
		&PackageManifest{
			PackageName:        "etcd",
			DefaultChannelName: "alpha",
			Channels: []PackageChannel{
				PackageChannel{Name: "alpha", CurrentCSVName: "etcdoperator.v0.9.2"},
				PackageChannel{Name: "nightly", CurrentCSVName: "etcdoperator.v0.9.2-nightly"},
			},
		},
	}

	// Do not use tabs for indentation as yaml forbids tabs http://yaml.org/faq.html
	rawCSV = `
data:
  clusterServiceVersions: |-
    - apiVersion: app.coreos.com/v1alpha1
      kind: ClusterServiceVersion-v1
      metadata:
        name: jbossapp-operator.v0.1.0
      spec:
        replaces: foo
        customresourcedefinitions:
          owned:
          - name: bar
            version: v1
            kind: JBossApp
          required:
          - name: baz
            version: v1
            kind: BazApp
`

	csvWant = []*olm_v1alpha1.ClusterServiceVersion{
		&olm_v1alpha1.ClusterServiceVersion{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "app.coreos.com/v1alpha1",
				Kind:       "ClusterServiceVersion-v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "jbossapp-operator.v0.1.0",
			},
			Spec: olm_v1alpha1.ClusterServiceVersionSpec{
				Replaces: "foo",
				CustomResourceDefinitions: olm_v1alpha1.CustomResourceDefinitions{
					Owned: []olm_v1alpha1.CRDDescription{
						olm_v1alpha1.CRDDescription{
							Name:    "bar",
							Version: "v1",
							Kind:    "JBossApp",
						},
					},
					Required: []olm_v1alpha1.CRDDescription{
						olm_v1alpha1.CRDDescription{
							Name:    "baz",
							Version: "v1",
							Kind:    "BazApp",
						},
					},
				},
			},
		},
	}
)

// Scenario: An operator manifest has a list of CRD(s).
// Expected Result: The list of CRD(s) gets marshaled into structured type.
func TestUnmarshal_ManifestHasCRD_SuccessfullyParsed(t *testing.T) {

	u := manifestYAMLParser{}
	dataGot, err := u.Unmarshal([]byte(rawCRDs))

	assert.NoError(t, err)
	assert.NotNil(t, dataGot)

	crdGot := dataGot.CustomResourceDefinitions

	assert.ElementsMatch(t, crdWant, crdGot)
}

// Scenario: An operator manifest has a list of package(s).
// Expected Result: The list of package(s) gets marshaled into structured type.
func TestUnmarshal_ManifestHasPackages_SuccessfullyParsed(t *testing.T) {
	u := manifestYAMLParser{}
	dataGot, err := u.Unmarshal([]byte(rawPackages))

	assert.NoError(t, err)
	assert.NotNil(t, dataGot)

	packagesGot := dataGot.Packages

	assert.ElementsMatch(t, packagesWant, packagesGot)
}

// Scenario: An operator manifest has a list of package(s).
// Expected Result: The list of package(s) gets marshaled into structured type.
func TestUnmarshal_ManifestHasCSV_SuccessfullyParsed(t *testing.T) {
	u := manifestYAMLParser{}
	dataGot, err := u.Unmarshal([]byte(rawCSV))

	assert.NoError(t, err)
	assert.NotNil(t, dataGot)

	csvGot := dataGot.ClusterServiceVersions

	assert.ElementsMatch(t, csvWant, csvGot)
}

// Given a structured representation of operator manifest we should be able to
// convert it to raw YAML representation so that a configMap object for catalog
// source can be created successfully.
func TestMarshal(t *testing.T) {
	marshaled := Manifest{
		CustomResourceDefinitions: crdWant,
		ClusterServiceVersions:    csvWant,
		Packages:                  packagesWant,
	}

	u := manifestYAMLParser{}
	rawGot, err := u.Marshal(&marshaled)

	assert.NoError(t, err)
	assert.NotNil(t, rawGot)
	assert.NotEmpty(t, rawGot.Packages)
	assert.NotEmpty(t, rawGot.CustomResourceDefinitions)
	assert.NotEmpty(t, rawGot.ClusterServiceVersions)
}

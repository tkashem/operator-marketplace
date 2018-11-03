package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	// Do not use tabs for indentation as yaml forbids tabs http://yaml.org/faq.html
	raw := `
publisher: redhat
data:
  customResourceDefinitions: "my crds"
  clusterServiceVersions: "my csvs"
  packages: "my packages"
`

	u := blobUnmarshalerImpl{}
	manifest, err := u.Unmarshal([]byte(raw))

	require.NoError(t, err)

	assert.Equal(t, "my crds", manifest.Data.CustomResourceDefinitions)
	assert.Equal(t, "my csvs", manifest.Data.ClusterServiceVersions)
	assert.Equal(t, "my packages", manifest.Data.Packages)
}

func TestRedHatOperatorManifest(t *testing.T) {
	raw := RedHatOperatorManifest

	u := blobUnmarshalerImpl{}
	manifest, err := u.Unmarshal([]byte(raw))
	require.NoError(t, err)

	pkg, err := u.UnmarshalData(&manifest.Data)
	require.NoError(t, err)

	m, err := u.ToManifest(pkg)
	assert.NotNil(t, m)
}

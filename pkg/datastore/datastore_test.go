package datastore

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPackageIDs(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	expected := []string{"foo/1", "bar/2", "braz/3"}

	packages := []*OperatorMetadata{
		helperNewOperatorMetadata("foo", "1"),
		helperNewOperatorMetadata("bar", "2"),
		helperNewOperatorMetadata("braz", "3"),
	}

	parser := NewMockManifestYAMLParser(controller)

	ds := &memoryDatastore{
		manifests: map[string]*OperatorManifest{},
		packages:  map[string]*ManifestPackage{},

		parser:   parser,
		unpacker: &manifestUnpacker{},
	}

	// We expect Unmarshal function to be invoked for each package.
	parser.EXPECT().Unmarshal(gomock.Any()).Return(&StructuredOperatorManifestData{}, nil).Times(len(packages))

	err := ds.Write(packages)
	require.NoError(t, err)

	result := ds.GetPackageIDs()
	actual := strings.Split(result, ",")

	assert.ElementsMatch(t, expected, actual)
}

func TestSomething(t *testing.T) {
	ds := New()

	rawYAML, err := ioutil.ReadFile("/home/akashem/rh-operators.yaml")
	require.NoError(t, err)

	packages := []*OperatorMetadata{
		&OperatorMetadata{
			RegistryMetadata: RegistryMetadata{
				Namespace:  "foo",
				Repository: "bar",
			},
			RawYAML: rawYAML,
		},
	}

	err = ds.Write(packages)
	assert.NoError(t, err)
}

func helperNewOperatorMetadata(namespace, repository string) *OperatorMetadata {
	return &OperatorMetadata{
		RegistryMetadata: RegistryMetadata{
			Namespace:  namespace,
			Repository: repository,
		},
	}
}

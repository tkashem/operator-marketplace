package datastore_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	"github.com/operator-framework/operator-marketplace/pkg/appregistry"
	"github.com/operator-framework/operator-marketplace/pkg/datastore"
)

func TestGetPackageIDs(t *testing.T) {
	packagesWant1 := []string{"foo/1", "bar/2", "braz/3"}
	packagesWant2 := []string{"a/v1", "b/v2", "c/v3"}

	opsrcID1, opsrcID2 := types.UID("123456"), types.UID("987654")

	opsrc1 := &v1alpha1.OperatorSource{
		ObjectMeta: metav1.ObjectMeta{
			UID: opsrcID1,
		},
	}
	opsrc2 := &v1alpha1.OperatorSource{
		ObjectMeta: metav1.ObjectMeta{
			UID: opsrcID2,
		},
	}

	manifests1 := []*appregistry.OperatorMetadata{
		&appregistry.OperatorMetadata{Namespace: "foo", Repository: "1"},
		&appregistry.OperatorMetadata{Namespace: "bar", Repository: "2"},
		&appregistry.OperatorMetadata{Namespace: "braz", Repository: "3"},
	}
	manifests2 := []*appregistry.OperatorMetadata{
		&appregistry.OperatorMetadata{Namespace: "a", Repository: "v1"},
		&appregistry.OperatorMetadata{Namespace: "b", Repository: "v2"},
		&appregistry.OperatorMetadata{Namespace: "c", Repository: "v3"},
	}

	ds := datastore.New()

	err := ds.Write(opsrc1, manifests1)
	require.NoError(t, err)
	err = ds.Write(opsrc2, manifests2)
	require.NoError(t, err)

	packages1 := ds.GetPackageIDs(opsrcID1)
	packages2 := ds.GetPackageIDs(opsrcID2)

	packagesGot1 := strings.Split(packages1, ",")
	packagesGot2 := strings.Split(packages2, ",")

	assert.ElementsMatch(t, packagesWant1, packagesGot1)
	assert.ElementsMatch(t, packagesWant2, packagesGot2)
}

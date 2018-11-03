package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	// Do not use tabs for indentation as yaml forbids tabs http://yaml.org/faq.html
	raw := `
data:
  customResourceDefinitions:
  clusterServiceVersions:
  packages:
`

	u := blobUnmarshalerImpl{}
	data, err := u.Unmarshal([]byte(raw))

	assert.NoError(t, err)
	assert.NotNil(t, data)
}

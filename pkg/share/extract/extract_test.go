package extract_test

import (
	"encoding/json"
	"testing"

	"github.com/http-everything/httpe/pkg/share/extract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractsFromI(t *testing.T) {
	var data interface{}
	JSON := []byte(`{
	 "name":"John",
	 "address": {
        "city":"London",
        "country":"UK"
     }}`)
	err := json.Unmarshal(JSON, &data)
	require.NoError(t, err)

	assert.Equal(t, "John", extract.SFromI("name", data))
	assert.Equal(t, "London", extract.SFromI("address.city", data))
	assert.Equal(t, "UK", extract.SFromI("address.country", data))
	assert.Equal(t, "", extract.SFromI("non.existing", data))
}

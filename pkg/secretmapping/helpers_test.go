package secretmapping_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jenkins-x/jx-secret/pkg/secretmapping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretMappingFind(t *testing.T) {
	dir := filepath.Join("test_data", "config-root", "namespaces")
	mapping, fileName, err := secretmapping.LoadSecretMapping(dir, true)
	require.NoError(t, err)
	require.NotEmpty(t, fileName, "no fileName returned")
	require.NotNil(t, mapping, "mapping not returned")

	testCases := []struct {
		secretName       string
		dataKey          string
		found            bool
		expectedKey      string
		expectedProperty string
	}{
		{
			secretName:       "lighthouse-hmac-token",
			dataKey:          "hmac",
			found:            true,
			expectedKey:      "secret/data/lighthouse/hmac",
			expectedProperty: "token",
		},
		{
			secretName:       "lighthouse-oauth-token",
			dataKey:          "oauth",
			found:            true,
			expectedKey:      "secret/data/jx/pipelineUser",
			expectedProperty: "token",
		},
	}

	for _, tc := range testCases {
		secretName := tc.secretName
		m := mapping.Find(secretName, tc.dataKey)
		if tc.found {
			require.NotNil(t, m, "should have found Mapping for secret %s and entry %s", secretName, tc.dataKey)

			assert.Equal(t, tc.expectedKey, m.Key, "key for secret %s", secretName)
			assert.Equal(t, tc.expectedProperty, m.Property, "property for secret %s", secretName)

			t.Logf("secret %s maps to key: %s property: %s\n", secretName, m.Key, m.Property)
		} else {
			assert.Nil(t, m, "should not have found Mapping for secret %s", secretName)
		}
	}
}

func TestMissingDefaultBackendType(t *testing.T) {
	sourceData := filepath.Join("test_data", "missing_default_backend_type")
	_, _, err := secretmapping.LoadSecretMapping(sourceData, true)
	require.Error(t, err, "failed did not receive error validating missing backend type")
	assert.True(t, strings.Contains(err.Error(), "Spec.Defaults.BackendType: zero value"))
}

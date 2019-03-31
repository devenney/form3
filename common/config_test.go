package common

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

// TestDefaultConfig tests that the configuration defaults are as expected.
func TestDefaultConfig(t *testing.T) {
	InitConfig()

	if viper.GetString(DBEndpointKey) != dbEndpointDefault {
		t.Fatalf("db_endpoint default was wrong. Expected: %s. Got: %s.",
			dbEndpointDefault, viper.GetString(DBEndpointKey))
	}
}

// TestEnvironmentOverrides tests that config values can be overridden
// by environment variables.
func TestEnvironmentOverrides(t *testing.T) {
	testCases := []struct {
		val string // The test value to set.
		key string // The key we are testing.
	}{
		{
			val: "url",
			key: APIUrlKey,
		},
		{
			val: "BIND_TEST:1234",
			key: BindAddressKey,
		},
		{
			val: "DBE_TEST",
			key: DBEndpointKey,
		},
		{
			val: "ENV_TEST",
			key: EnvKey,
		},
	}

	for _, test := range testCases {
		os.Setenv(fmt.Sprintf("%s_%s", envPrefix, test.key), test.val)

		InitConfig()

		actual := viper.GetString(test.key)
		if actual != test.val {
			t.Fatalf("%s value was wrong. Expected: %s. Got: %s", test.key, test.val, actual)
		}
	}
}

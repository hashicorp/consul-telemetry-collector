package confresolver

import (
	"testing"

	"github.com/shoenig/test"
)

func TestComponentConfig(t *testing.T) {
	ccfg := make(componentConfig)
	const keyValue = "key-value"
	const otherKeyValue = "other-key-value"
	const secretValue = "secret-value"
	ccfg.Set("key", keyValue)
	ccfg.Set("otherkey", otherKeyValue)
	ccfg.Set("secret", secretValue)
	ccfg.SetMap("map").Set("key", keyValue)

	expectPresentAndValue(t, ccfg, "key", keyValue)
	expectPresentAndValue(t, ccfg, "otherkey", otherKeyValue)
	expectPresentAndValue(t, ccfg, "secret", secretValue)

	containerInterface, ok := ccfg["map"]
	test.True(t, ok)
	innerCcfg, ok := containerInterface.(componentConfig)
	test.True(t, ok)
	expectPresentAndValue(t, innerCcfg, "key", keyValue)
}

func expectPresentAndValue(t *testing.T, container componentConfig, key string, expectedValue interface{}) {
	val, ok := container[key]
	test.True(t, ok)
	test.Eq(t, expectedValue, val)
}

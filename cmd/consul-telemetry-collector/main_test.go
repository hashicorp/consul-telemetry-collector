package main

import (
	"flag"
	"os"
	"testing"

	"github.com/shoenig/test"
)

func Test_envVarString(t *testing.T) {
	const flagName = "name"
	const envKey = "TEST_NAME_KEY"
	const defaultValue = "This is the unset value"
	testcases := map[string]struct {
		flags  []string
		env    map[string]string
		expect string
	}{
		"EnvOnly": {
			env:    map[string]string{envKey: "env-only"},
			expect: "env-only",
		},
		"FlagOnly": {
			flags:  []string{"-name", "flag-only"},
			expect: "flag-only",
		},
		"FlagOverrideEnv": {
			env:    map[string]string{envKey: "env-value"},
			flags:  []string{"-name", "flag-value"},
			expect: "flag-value",
		},
		"NoValue": {
			expect: defaultValue,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			// reset flag.CommandLine since we call parse multiple times
			fs := flag.NewFlagSet(name, flag.ExitOnError)
			flag.CommandLine = fs
			setupEnv(t, tc.env)
			var ptrString string
			flag.StringVar(&ptrString, flagName, defaultValue, "usage")
			envVarString(envKey, &ptrString)
			test.NoError(t, fs.Parse(tc.flags))
			test.Eq(t, ptrString, tc.expect)
		})
	}
}

func setupEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for k, v := range env {
		os.Setenv(k, v)
	}
	t.Cleanup(func() {
		for k := range env {
			os.Unsetenv(k)
		}
	})
}

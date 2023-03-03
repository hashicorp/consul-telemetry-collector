package main

import (
	"flag"
	"fmt"
	"os"
)

func stringVarOrEnv(flagValue *string, name, defaultValue, usage, envKey string) {
	usage = appendEnvUsage(envKey, usage)
	flag.StringVar(flagValue, name, defaultValue, usage)

	envVal, ok := os.LookupEnv(envKey)
	if ok {
		// we found this in the environment
		*flagValue = envVal
		return
	}
}

func appendEnvUsage(envKey string, usage string) string {
	return fmt.Sprintf("%s\n\tEnvironment variable %s", usage, envKey)
}

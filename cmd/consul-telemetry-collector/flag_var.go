package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func stringVar(ptr *string, name, defaultValue, usage, envKey string) {
	usage = appendEnvUsage(envKey, usage)
	flag.StringVar(ptr, name, defaultValue, usage)
	*ptr = must(parseEnv(envKey, defaultValue, func(s string) (string, error) {
		return s, nil
	}))
}

func parseEnv[T any](envKey string, defaultValue T, parseFn func(string) (T, error)) (T, error) {
	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		return defaultValue, nil
	}
	val, err := parseFn(envVal)
	if err != nil {
		return defaultValue, fmt.Errorf("unable to parse environment variable %s=%s as %T", envKey, envVal, val)
	}
	return val, nil
}

func appendEnvUsage(envKey string, usage string) string {
	return fmt.Sprintf("%s\n\tEnvironment variable %s", usage, envKey)
}

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return v
}

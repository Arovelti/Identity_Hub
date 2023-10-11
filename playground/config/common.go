package config

import (
	"log"
	"os"
	"strconv"
)

func MustString(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("required ENV %q is not set", key)
	}
	if value == "" {
		log.Fatalf("required ENV %q is empty", key)
	}

	return value
}

func MustInt(key string) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("required ENV %s is not set", key)
	}
	if value == "" {
		log.Fatalf("required ENV %q is empty", key)
	}
	res, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		log.Fatalf("required ENV %s must be a number but it's %s", key, value)
	}

	return int(res)
}

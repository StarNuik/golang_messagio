package lib

import "os"

func EnvOr(key string, or string) string {
	val := os.Getenv(key)
	if val == "" {
		val = or
	}
	return val
}

package helpers

import "os"

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetEnvOrDefault(key string, def string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	} else {
		return def
	}
}

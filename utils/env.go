package utils

import (
	"os"
	"strings"
)

func EnvMap() map[string]string {
	environment := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		environment[pair[0]] = pair[1]
	}
	return environment
}

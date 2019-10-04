package utils

import (
	"os"
	"strings"
)

func EnvMap() map[string]string {
	envrionment := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envrionment[pair[0]] = pair[1]
	}
	return envrionment
}

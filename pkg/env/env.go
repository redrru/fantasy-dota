package env

import (
	"os"
	"strings"
	"sync"
)

type Env map[string]string

var (
	env  = Env{}
	once sync.Once
)

func GetEnv() Env {
	once.Do(func() {
		for _, value := range os.Environ() {
			split := strings.Split(value, "=")
			if len(split) != 2 {
				continue
			}
			env[split[0]] = split[1]
		}
	})

	return env
}

package env

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/redrru/fantasy-dota/pkg/log"
)

type Env struct {
	store map[string]string
}

func (e *Env) GetString(k string) string {
	return e.store[k]
}

func (e *Env) GetInt(k string) int {
	v, _ := strconv.ParseInt(e.store[k], 10, 64)
	return int(v)
}

func (e *Env) GetDuration(k string) time.Duration {
	v, _ := time.ParseDuration(e.store[k])
	return v
}

var (
	env  = Env{store: map[string]string{}}
	once sync.Once
)

func GetEnv() Env {
	once.Do(func() {
		for _, value := range os.Environ() {
			split := strings.SplitN(value, "=", 2)
			if len(split) != 2 {
				continue
			}
			env.store[split[0]] = split[1]
		}

		log.GetLogger().Info(context.Background(), "[Env] GetEnv", zap.Any("store", env.store))
	})

	return env
}

//go:build unit
// +build unit

package env

import (
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	os.Clearenv()

	envs := Env{}
	for i := 0; i < 1000; i++ {
		k, v := gofakeit.Word(), gofakeit.Word()
		err := os.Setenv(k, v)
		assert.NoError(t, err)
		envs[k] = v
	}

	result := GetEnv()
	assert.Equal(t, envs, result)
}

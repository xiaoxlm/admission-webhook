package pkg

import (
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func Test_checkEnv(t *testing.T) {
	err := checkEnv([]corev1.EnvVar{
		{
			Name:  MustKey,
			Value: "111",
		},
	})

	assert.Equal(t, nil, err)
}

package pkg

import (
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestInjectMappingData_Mutating(t *testing.T) {
	pod := &corev1.Pod{}
	inject := &InjectMappingData{
		raw: pod,
	}

	_ = inject.Mutating()

	_, ok := inject.mutated.GetAnnotations()["mutate-timestamp"]
	assert.Equal(t, true, ok)
	assert.Equal(t, map[string]string{
		"mutated-app": "true",
	}, inject.mutated.GetLabels())

}

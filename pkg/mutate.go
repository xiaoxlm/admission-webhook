package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/wI2L/jsondiff"
	corev1 "k8s.io/api/core/v1"
	"log"
	"net/http"
	"time"
)

func Mutate(w http.ResponseWriter, r *http.Request) {
	reviewRequest, reviewRESP, err := InitReviewRequestResponse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = checkKind(reviewRequest.Request.Kind.Kind); err != nil {
		FailureReviewRESP(reviewRESP, err.Error())
		Response(reviewRESP, w)
		return
	}

	pod, err := GetPod(reviewRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("mutating pod: ", pod.Name)
	if pod.GetNamespace() != "webhook" {
		Response(reviewRESP, w)
		return
	}

	patchData, err := NewInjectMappingData(pod).Mutating().GetPatchData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	PatchReviewRESP(reviewRESP, patchData)
	Response(reviewRESP, w)
	return
}

type InjectMappingData struct {
	raw     *corev1.Pod
	mutated *corev1.Pod
}

func NewInjectMappingData(raw *corev1.Pod) *InjectMappingData {
	return &InjectMappingData{
		raw: raw,
	}
}

func (inject *InjectMappingData) Mutating() *InjectMappingData {
	mpod := inject.raw.DeepCopy()

	inject.injectAnnotation(mpod)
	inject.injectLabel(mpod)

	inject.mutated = mpod

	return inject
}

func (inject *InjectMappingData) GetPatchData() ([]byte, error) {
	patch, err := jsondiff.Compare(inject.raw, inject.mutated)
	if err != nil {
		return nil, err
	}

	return json.Marshal(patch)
}

func (inject *InjectMappingData) injectAnnotation(pod *corev1.Pod) {
	annotations := pod.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["mutate-timestamp"] = fmt.Sprintf("%d", time.Now().Unix())

	pod.SetAnnotations(annotations)
}

func (inject *InjectMappingData) injectLabel(pod *corev1.Pod) {
	labels := pod.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["mutated-app"] = "true"

	pod.SetLabels(labels)
}

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

	patchData, err := mutating(pod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	PatchReviewRESP(reviewRESP, patchData)
	Response(reviewRESP, w)
	return
}

func mutating(pod *corev1.Pod) ([]byte, error) {
	mpod := pod.DeepCopy()

	injectAnnotation(mpod)
	injectLabel(mpod)

	// generate json patch
	patch, err := jsondiff.Compare(pod, mpod)
	if err != nil {
		return nil, err
	}

	return json.Marshal(patch)
}

func injectAnnotation(pod *corev1.Pod) {
	annotations := pod.GetAnnotations()
	annotations["mutate-timestamp"] = fmt.Sprintf("%d", time.Now().Unix())

	pod.SetAnnotations(annotations)
}

func injectLabel(pod *corev1.Pod) {
	labels := pod.GetLabels()
	labels["mutated-app"] = "true"

	pod.SetLabels(labels)
}

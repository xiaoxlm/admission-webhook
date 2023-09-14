package pkg

import (
	"encoding/json"
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
	"net/http"
)

const (
	MustKey = "Time"
)

func GetPod(reviewRequest *admissionv1.AdmissionReview) (*corev1.Pod, error) {
	var pod = &corev1.Pod{}
	if err := json.Unmarshal(reviewRequest.Request.Object.Raw, pod); err != nil {
		return nil, fmt.Errorf("pod unmarshal error:" + err.Error())
		//http.Error(w, "pod unmarshal error:"+err.Error(), http.StatusInternalServerError)
		//return
	}

	return pod, nil
}

func Validate(w http.ResponseWriter, r *http.Request) {
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
	log.Println("validating pod: ", pod.Name)

	if pod.GetNamespace() != "webhook" {
		Response(reviewRESP, w)
		return
	}

	for _, container := range pod.Spec.Containers {
		if err = checkEnv(container.Env); err != nil {
			FailureReviewRESP(reviewRESP, fmt.Sprintf("container %s validate failed.%s", container.Name, err.Error()))
			break
		}
	}

	Response(reviewRESP, w)
	return
}

func checkEnv(envs []corev1.EnvVar) error {
	if len(envs) == 0 {
		return fmt.Errorf("env vars is empty")
	}

	{
		var withTimeKey bool
		for _, e := range envs {
			if e.Name == MustKey {
				withTimeKey = true
			}
		}
		if !withTimeKey {
			return fmt.Errorf(fmt.Sprintf("env vars doesn't have '%s' key", MustKey))
		}
	}

	return nil
}

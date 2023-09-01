package main

import (
	"encoding/json"
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	TLSCertName = "tls.crt"
	TLSKeyName  = "tls.key"
	MustKey     = "Time"
)

func main() {
	var (
		certFile string
		keyFile  string
	)
	if certDir := os.Getenv("CERT_DIR"); certDir != "" {
		certFile = filepath.Join(certDir, TLSCertName)
		keyFile = filepath.Join(certDir, TLSKeyName)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validate)

	log.Println("server listen on ':8000'...")

	if certFile != "" && keyFile != "" {
		log.Fatal(http.ListenAndServeTLS(":8000", certFile, keyFile, mux))
	} else {
		log.Fatal(http.ListenAndServe(":8000", mux))
	}
}

func validate(w http.ResponseWriter, r *http.Request) {
	var reviewRequest admissionv1.AdmissionReview
	{
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&reviewRequest); err != nil {
			http.Error(w, "admission review request body decode error:"+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var pod = corev1.Pod{}
	if err := json.Unmarshal(reviewRequest.Request.Object.Raw, &pod); err != nil {
		http.Error(w, "pod unmarshal error:"+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("validated pod ", pod.Name)

	var reviewRESP = admissionv1.AdmissionReview{
		TypeMeta: reviewRequest.TypeMeta,
		Response: &admissionv1.AdmissionResponse{
			UID:     reviewRequest.Request.UID, // write the unique identifier back
			Allowed: true,
		},
	}

	for _, container := range pod.Spec.Containers {
		err := checkEnv(container.Env)

		if err != nil {
			reviewRESP.Response.Allowed = false
			reviewRESP.Response.Result = &metav1.Status{
				Status:  "Failure",
				Message: "pod validating invalid",
				Reason:  metav1.StatusReason(fmt.Sprintf("container %s validate failed.%s", container.Name, err.Error())),
				Code:    402,
			}
			break
		}
	}

	resp, _ := json.Marshal(reviewRESP)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
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

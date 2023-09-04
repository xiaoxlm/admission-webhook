package pkg

import (
	"encoding/json"
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func InitReviewRequestResponse(r *http.Request) (reviewRequest, reviewRESP *admissionv1.AdmissionReview, err error) {
	reviewRequest = &admissionv1.AdmissionReview{}
	{
		dec := json.NewDecoder(r.Body)
		if err = dec.Decode(&reviewRequest); err != nil {
			err = fmt.Errorf("admission review request body decode error:" + err.Error())
			return
		}
	}

	reviewRESP = &admissionv1.AdmissionReview{
		TypeMeta: reviewRequest.TypeMeta,
		Response: &admissionv1.AdmissionResponse{
			UID:     reviewRequest.Request.UID, // write the unique identifier back
			Allowed: true,
		},
	}

	return
}

func Response(reviewRESP *admissionv1.AdmissionReview, w http.ResponseWriter) {
	resp, _ := json.Marshal(reviewRESP)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	return
}

func FailureReviewRESP(reviewRESP *admissionv1.AdmissionReview, reason string) {
	reviewRESP.Response.Allowed = false
	reviewRESP.Response.Result = &metav1.Status{
		Status:  "Failure",
		Message: "pod validating invalid",
		Reason:  metav1.StatusReason(reason),
		Code:    http.StatusBadRequest,
	}

	return
}

func PatchReviewRESP(reviewRESP *admissionv1.AdmissionReview, patchData []byte) {
	patchType := admissionv1.PatchTypeJSONPatch

	reviewRESP.Response.PatchType = &patchType
	reviewRESP.Response.Patch = patchData

	return
}

func checkKind(kind string) error {
	if kind != "Pod" {
		return fmt.Errorf("the resource kind is not pod")
	}

	return nil
}

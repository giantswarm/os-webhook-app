package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
)

type AdmissionHookInterface interface {
	//Return the admission review
	GetAdmissionReview(r *http.Request) *v1beta1.AdmissionReview

	//Returns answer to the API
	Answer(w http.ResponseWriter, r *http.Request, answer *admissionv1beta1.AdmissionResponse)
}

type AdmissionWebhook struct {
	AdmissionHookInterface
}

func (ah *AdmissionWebhook) GetAdmissionReview(w http.ResponseWriter, r *http.Request) *v1beta1.AdmissionReview {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		glog.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return nil
	}
	glog.Info("Received request")

	arRequest := new(v1beta1.AdmissionReview)
	if err := json.Unmarshal(body, &arRequest); err != nil {
		glog.Error("incorrect body")
		http.Error(w, "incorrect body", http.StatusBadRequest)
	}

	return arRequest
}

func (ah *AdmissionWebhook) Answer(w http.ResponseWriter, r *http.Request, answer *admissionv1beta1.AdmissionResponse) error {
	resp, err := json.Marshal(answer)
	if err != nil {
		glog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return err
	}

	glog.Infof("Ready to write reponse ...")
	if _, err := w.Write(resp); err != nil {
		glog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
		return err
	}

	return nil
}

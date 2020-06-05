package webhook

import (
	"fmt"
	"net/http"
	"os"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/golang/glog"
	log "github.com/sirupsen/logrus"
)

type AWSClusterAdmissionHook struct {
	AdmissionWebhook

	client *kubernetes.Clientset // Kubernetes client for calling Api
}

func (ah *AWSClusterAdmissionHook) Admit(req *admissionv1beta1.AdmissionRequest) *admissionv1beta1.AdmissionResponse {
	resp := &admissionv1beta1.AdmissionResponse{}
	resp.UID = req.UID

	// Skip operations that aren't create or update
	if req.Operation != admissionv1beta1.Create &&
		req.Operation != admissionv1beta1.Update {
		glog.Infof("Skipping %s request for %s", req.Operation, req.Object)
		resp.Allowed = true
		return resp
	}

	glog.Infof("Processing %s request for %s", req.Operation, req.Name)
	glog.V(2).Infof("Incoming object: %s", req.Object)

	resp.Allowed = true
	return resp
}

func (ah *AWSClusterAdmissionHook) Initialize(kubeClientConfig *rest.Config) error {
	// Initialise a Kubernetes client
	client, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return fmt.Errorf("failed to intialise kubernetes clientset: %v", err)
	}
	ah.client = client

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)

	log.Info("Webhook Initialization Complete.")
	return nil
}

func (ah *AWSClusterAdmissionHook) Serve(w http.ResponseWriter, r *http.Request) {
	admissionReview := ah.GetAdmissionReview(w, r)
	glog.V(5).Infof("Incoming Admission review: %+v\n", admissionReview)
	admissionResponse := ah.Admit(admissionReview.Request)
	ah.Answer(w, r, admissionResponse)
}

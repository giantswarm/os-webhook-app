package webhook

import (
	"fmt"
	"os"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	log "github.com/sirupsen/logrus"
)

// AdmissionHook implements the OpenShift MutatingAdmissionHook interface.
// https://github.com/openshift/generic-admission-server/blob/v1.9.0/pkg/apiserver/apiserver.go#L45
type AWSClusterAdmissionHook struct {
	client *kubernetes.Clientset // Kubernetes client for calling Api
}

func (ah *AWSClusterAdmissionHook) ValidatingResource() (plural schema.GroupVersionResource, singular string) {
	return schema.GroupVersionResource{
			Group:    "infrastructure.giantswarm.io",
			Version:  "v1alpha2",
			Resource: "awsclusters",
		},
		"AdmissionReview"
}

func (ah *AWSClusterAdmissionHook) Admit(req *admissionv1beta1.AdmissionRequest) *admissionv1beta1.AdmissionResponse {
	resp := &admissionv1beta1.AdmissionResponse{}
	resp.UID = req.UID
	requestName := fmt.Sprintf("%s %s", req.Kind, req.Object)

	// Skip operations that aren't create or update
	if req.Operation != admissionv1beta1.Create &&
		req.Operation != admissionv1beta1.Update {
		log.Info("Skipping %s request for %s", req.Operation, requestName)
		resp.Allowed = true
		return resp
	}

	log.Info("Processing %s request for %s", req.Operation, requestName)

	resp.Allowed = true
	return resp
}

func (ah *AWSClusterAdmissionHook) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
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

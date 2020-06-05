package main

import (
	"github.com/giantswarm/os-webhook-app/webhook"
	"github.com/openshift/generic-admission-server/pkg/cmd"
)

func main() {
	cmd.RunAdmissionServer(
		&webhook.AWSControlPlaneAdmissionHook{},
		&webhook.AWSClusterAdmissionHook{},
	)
}

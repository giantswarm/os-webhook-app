module github.com/giantswarm/os-webhook-app

go 1.14

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/openshift/generic-admission-server v1.14.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.0.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/apiserver v0.18.2 // indirect
	k8s.io/client-go v0.18.2
	k8s.io/utils v0.0.0-20200603063816-c1c6865ac451 // indirect
)

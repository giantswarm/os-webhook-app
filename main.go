package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/giantswarm/os-webhook-app/webhook"
	"github.com/golang/glog"
	restclient "k8s.io/client-go/rest"
)

const (
	port = "4430"
)

var (
	tlscert, tlskey string
)

func main() {

	flag.StringVar(&tlscert, "tlsCertFile", "/etc/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&tlskey, "tlsKeyFile", "/etc/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")

	flag.Parse()

	certs, err := tls.LoadX509KeyPair(tlscert, tlskey)
	if err != nil {
		glog.Errorf("Filed to load key pair: %v", err)
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%v", port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{certs}},
	}

	inClusterConfig, err := restclient.InClusterConfig()
	if err != nil {
		glog.Errorf("Filed to load key kubeconfig: %v", err)
	}
	// define http server and server handler
	awscluster := &webhook.AWSClusterAdmissionHook{}
	awscluster.Initialize(inClusterConfig)

	mux := http.NewServeMux()
	mux.HandleFunc("/awscluster", awscluster.Serve)
	server.Handler = mux

	glog.Infof("Endpoints initialized")

	// start webhook server in new rountine
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			glog.Errorf("Failed to listen and serve webhook server: %v", err)
		}
	}()

	glog.Infof("Server running listening in port: %s", port)

	// listening shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	glog.Info("Got shutdown signal, shutting down webhook server gracefully...")
	server.Shutdown(context.Background())
}

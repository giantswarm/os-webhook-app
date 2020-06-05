CGO_ENABLED=0 go build .
docker build . -t paurosello/os-webhook-app:0.0.17
docker push paurosello/os-webhook-app:0.0.17

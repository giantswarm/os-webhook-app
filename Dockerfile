FROM centos
WORKDIR /app
COPY os-webhook-app /app
CMD ["/app/os-webhook-app"]

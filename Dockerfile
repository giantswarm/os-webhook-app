FROM alpine:3.11
WORKDIR /app
COPY os-webhook-app /app
CMD ["/app/os-webhook-app"]

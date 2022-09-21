# K8s

## Deploy services

* `downstream-service.yml` - dowtream service - dom's implementation of unstable server
* `echo-service` - ehco service written in go, returns plain string. Used for testing HPA with custom metrics
* `go-service.yml` - GoLang implementation of server
* `py-service.yml` - Python implementation of server
* `py-ev-service.yml` - another Python implementation of server
* `ts-service.yml` - TypeScript implementation of server
* `tools.yml` - handy toools used for loadtesting
* `istio-ingress-gateway.yml`  - contains istio ingress gateway, so the services are accessible when using istio
* `prometheus-adapter-values.yml` - yaml values for promoetheus adapter Heml release, containing definition of custom metrics for hpa

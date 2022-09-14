# Decisions

## Service mesh

Use istio to provide

* provide visibility over requests coming to the service
  * based on custom metrics HPA spins up more pods
* rate limits to prevent service to be overhelmed
  * https://istio.io/v1.12/docs/tasks/policy-enforcement/rate-limit/

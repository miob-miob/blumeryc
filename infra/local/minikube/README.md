https://istio.io/latest/docs/setup/platform-setup/minikube/

# Minikube

Start minikube

```
minikube config set driver docker
minikube start --memory=10240 --cpus=4 --kubernetes-version=v1.24.3
```

```
https://medium.com/codex/setup-istio-ingress-traffic-management-on-minikube-725c5e6d767a
```


## Deploy istio service mesh

based on https://www.arthurkoziel.com/running-knative-with-istio-in-kind/

get `istioctl` https://istio.io/latest/docs/reference/commands/istioctl/

```
kubectl label namespace default istio-injection=enabled
istioctl manifest apply --set profile=default \
--set components.egressGateways[0].name=istio-egressgateway \
--set components.egressGateways[0].enabled=true

# Access to external service \
# --set meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY \

```

## Kiali

```
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.15/samples/addons/kiali.yaml
```

## Prometheus



kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.15/samples/addons/prometheus.yaml

```
istioctl dashboard prometheus
```

```
sum(rate(istio_requests_total{app="echo", connection_security_policy="mutual_tls"}[1m])) by (destination_workload, code)
```

## Helm

https://helm.sh/docs/intro/install/

## Prometheus adapter

```
# https://artifacthub.io/packages/helm/prometheus-community/prometheus-adapter
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade prometheus-adapter prometheus-community/prometheus-adapter --install --create-namespace --namespace monitoring-system -f prometheus-adapter-values.yml
```


## Set rate limit


https://dev.to/tresmonauten/setup-an-ingress-rate-limiter-with-envoy-and-istio-1i9g

## HPA

https://github.com/stefanprodan/istio-hpa

https://martinheinz.dev/blog/76
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/istio_requests_per_second"  | jq .
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/services/*/istio_requests_per_second"  | jq .


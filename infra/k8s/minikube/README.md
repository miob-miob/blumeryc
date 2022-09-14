# Minikube

get `minikube` https://minikube.sigs.k8s.io/docs/start/

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

kubectl apply -n istio-system -f istio-ingress-gateway.yml 
```


## Prometheus


kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.15/samples/addons/prometheus.yaml

```
istioctl dashboard prometheus
```

## Helm

get `helm` https://helm.sh/docs/intro/install/


## Prometheus adapter

```
# https://artifacthub.io/packages/helm/prometheus-community/prometheus-adapter
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade prometheus-adapter prometheus-community/prometheus-adapter --install --create-namespace --namespace monitoring-system -f prometheus-adapter-values.yml
```

## HPA

kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/istio_requests_per_second"  | jq .
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/services/*/istio_requests_per_second"  | jq .

https://github.com/stefanprodan/istio-hpa
https://martinheinz.dev/blog/76
s


## Kiali

```
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.15/samples/addons/kiali.yaml
```

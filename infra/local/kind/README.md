# Local k8s kind 

Stack

* kind https://kind.sigs.k8s.io/docs/user/quick-start/
* nginx ingress controller
 * based on https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

## Deploy cluster

Spin ip `kind` cluster

```
kind  create cluster --config=kind.yml
```

Deploy `nginx ingress controller`

```
kubectl apply -f ingress-controller.yml
```

Deploy `echo service`

```
kubectl apply -f echo-service.yml
curl http://localhost:8080/
```


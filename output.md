
# Output
 
![logo](./logo.png)

## how the selected technologies supported parallelism (and concurrency),

- GO
  - props
    - compiled language (fast, memory efficient)
    - language is designed to do async stuffs out of the box
  - cons
    - harder syntax to handle parallelism code	 
- TS
  - props
    - thanks to the event loop support of async code out of the box
    - rich tooling to handle those edge cases
  - cons
    - main runtime CPU computing is on the single system process (could be slower)
    - transpiled language
- PY
  - props:
    - thanks to [asyncio](https://docs.python.org/3/library/asyncio.html) (asgi) we have similar approach like the node js has
  - cons
    - hard to write async code with the most streight forward way	(HTTP server with threads per request in apache + mod_wsgi)
    - transpiled language
  - **fun facts:** 
    - For this kind of assignement we picked python with event loop to not to waste resources with creating 
      new threds which will wait till thre reponse from downstream service will come.
  
conclusion: 

it really depends on what does __high__ performance really means. If we can afford 
node, development cost may be cheaper. Weakest point of application also matters  - in our case it appears to be downstream service 

In this case the pick of the technology does not mean so much as we expected at the start of the test,
because the bottleneck of this app was in the third part party downstream service.

If we'll have more CPU heavy loaded algorithms in our codebase, 
GO will have much better performance in comparision with interpreted languages.

Of course as you may see in the **load testing** section, GO is the fastest and safest implementation at all. 
The second largest language is the Typescript and the Python is the slowest one for our test cases.


--------------- 

## Please consider Kubernetes (GKE) 

- Considered

POC environemnt implemented in minikube environment, to test integration of all used components and services.

Deployment to GKE might happen to test on real resources with higher load and more real numbers

## What tools will you use for autoscaling? 

Using HPA, see what metrics we would use and how.

## Suggest the best way for resource allocation and how we may solve the autoscaling

Using HPA, see what metrics we would use and how.

## Will you prefer horizontal or vertical autoscaling?

Horizontal autoscaling, as we believe more concurrency on service instance level is more effective for most of implementations.

## What metrics will you use as inputs for autoscaling?

As we decided not to use vertical autoscaling, we would like to know how much requests is going to our service and based on those metrics we would scale new replicas. For having such metrics, we suggest to use some kind of service mesh (istio), with help of prometheus and prometheus-adapter can be cconvertes to custom metrics used by horizontal pod autoscaler in kubernetes. As baseline for metrics we could use numbers from load tests. 

## simple API performance tests which

Used `vegeta` load test tool https://github.com/tsenart/vegeta

### How will the different approaches in parallelism (and concurrency) affect the solution's scalability?

see tests results

## Others

### Preventing server against overload

Using HPA should take care of spinning more replicas under load. We can leverage the istio rate limits to project over server as well (https://istio.io/latest/docs/tasks/policy-enforcement/rate-limit/ - not tested).

------------------------------------------------------------------------
------------------------------------------------------------------------
------------------------------------------------------------------------

# Load testing

## All in cluster

Deployed to AKS cluster, without istio and without hpa based on custom metrics

* downstream service deployed to cluster
* each implemantion deployed to cluster
* load testing tooling deployed and tests executed inside the cluster

## Use cases

x:y:z

where x = number of downastream replicas (doms implementation)
where y = number of implementation replicas
where z = number of concurrent requests

other `vegeta` params = `-duration=1m -workers=100`

### 1:1:100

baseline:
```
echo "GET http://downstream:3333" | vegeta attack -rate=100  -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         6000, 100.02, 82.73
Duration      [total, attack, wait]             1m0s, 59.99s, 469.803ms
Latencies     [min, mean, 50, 90, 95, 99, max]  929.818??s, 299.196ms, 299.314ms, 537.645ms, 569.757ms, 594.31ms, 601.624ms
Bytes In      [total, mean]                     346890, 57.81
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           83.37%
Status Codes  [code:count]                      200:5002  400:166  408:184  500:323  502:156  503:169
Error Set:
408 Request Timeout
500 Internal Server Error
502 Bad Gateway
503 Service Unavailable
400 Bad Request
```

py:
```
echo "GET http://py:9002/py?timeout=333" | vegeta attack -rate=100 -duration=1m -workers=100 | 
Requests      [total, rate, throughput]         6000, 100.02, 1.77
Duration      [total, attack, wait]             1m30s, 59.99s, 30s
Latencies     [min, mean, 50, 90, 95, 99, max]  188.843ms, 29.562s, 30s, 30s, 30s, 30.001s, 30.014s
Bytes In      [total, mean]                     13968, 2.33
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           2.65%
Status Codes  [code:count]                      0:5673  200:159  422:168
Error Set:
422 UNPROCESSABLE ENTITY
...Number of `read: connection reset by peer`...
```

py-ev:
```
echo "GET http://py-ev:9002/py-ev?timeout=333" | vegeta attack -rate=100 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         6000, 100.02, 48.93
Duration      [total, attack, wait]             1m0s, 59.99s, 337.884ms
Latencies     [min, mean, 50, 90, 95, 99, max]  4.781ms, 258.552ms, 320.054ms, 376.943ms, 379.02ms, 382.964ms, 402.825ms
Bytes In      [total, mean]                     242292, 40.38
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           49.20%
Status Codes  [code:count]                      200:2952  422:3048
Error Set:
422 UNPROCESSABLE ENTITY
Get "http://py:9002/py?timeout=333": EOF
Get "http://py:9002/py?timeout=333": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

ts:
```
echo "GET http://ts:2020/ts?timeout=333" | vegeta attack -rate=100 -duration=1m -workers=100 | vegeta report

Requests      [total, rate, throughput]         6000, 100.02, 46.56
Duration      [total, attack, wait]             1m0s, 59.99s, 323.311ms
Latencies     [min, mean, 50, 90, 95, 99, max]  2.329ms, 244.614ms, 306.236ms, 334.902ms, 335.287ms, 337.559ms, 512.168ms
Bytes In      [total, mean]                     321297, 53.55
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           46.80%
Status Codes  [code:count]                      200:2808  422:3192
Error Set:
422 Unprocessable Entity
```

go:
```
echo "GET http://go:8090/go?timeout=333" | vegeta attack -rate=100 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         6000, 100.02, 48.82
Duration      [total, attack, wait]             1m0s, 59.99s, 334.563ms
Latencies     [min, mean, 50, 90, 95, 99, max]  1.695ms, 254.693ms, 306.391ms, 334.471ms, 334.774ms, 336.491ms, 343.041ms
Bytes In      [total, mean]                     324471, 54.08
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           49.08%
Status Codes  [code:count]                      200:2945  422:3055
Error Set:
422 Unprocessable Entity
```

### 1:1:500

baseline:
```
echo "GET http://downstream:3333" | vegeta attack -rate=500  -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         30000, 500.02, 410.30
Duration      [total, attack, wait]             1m1s, 59.998s, 567.073ms
Latencies     [min, mean, 50, 90, 95, 99, max]  830.616??s, 301.117ms, 302.824ms, 541.034ms, 570.752ms, 594.826ms, 603.224ms
Bytes In      [total, mean]                     1730999, 57.70
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           82.83%
Status Codes  [code:count]                      200:24850  400:859  408:863  500:1730  502:845  503:853
Error Set:
400 Bad Request
500 Internal Server Error
408 Request Timeout
503 Service Unavailable
502 Bad Gateway
```

py:
```
echo "GET http://py:9002/py?timeout=333" | vegeta attack -rate=500 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         29995, 499.90, 1.73
Duration      [total, attack, wait]             1m30s, 1m0s, 29.92s
Latencies     [min, mean, 50, 90, 95, 99, max]  3.913ms, 28.483s, 30s, 30.001s, 30.003s, 30.035s, 30.486s
Bytes In      [total, mean]                     13560, 0.45
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           0.52%
Status Codes  [code:count]                      0:29680  200:156  422:159
Error Set:
422 UNPROCESSABLE ENTITY
...Number of `read: connection reset by peer`...
```

py-ev:
```
echo "GET http://py-ev:9002/py-ev?timeout=333" | vegeta attack -rate=500 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         30000, 500.02, 1.03
Duration      [total, attack, wait]             1m30s, 59.998s, 30s
Latencies     [min, mean, 50, 90, 95, 99, max]  29.282ms, 24.775s, 30s, 30s, 30s, 30.003s, 30.247s
Bytes In      [total, mean]                     182076, 6.07
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           0.31%
Status Codes  [code:count]                      0:18908  200:93  422:10999
Error Set:
422 Unprocessable Entity
...Number of `read: connection reset by peer`...
```

ts:
```
echo "GET http://ts:2020/ts?timeout=333" | vegeta attack -rate=500 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         30000, 500.02, 32.50
Duration      [total, attack, wait]             1m30s, 59.998s, 29.946s
Latencies     [min, mean, 50, 90, 95, 99, max]  3.036ms, 20.248s, 30s, 30s, 30.001s, 30.003s, 30.317s
Bytes In      [total, mean]                     598789, 19.96
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           9.74%
Status Codes  [code:count]                      0:17607  200:2923  422:9470
Error Set:
422 Unprocessable Entity
...Number of `Get "http://ts:2020/ts?timeout=333": context deadline exceeded (Client.Timeout exceeded while awaiting headers)`
```

go:
```
echo "GET http://go:8090/go?timeout=333" | vegeta attack -rate=500 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         30000, 500.00, 125.31
Duration      [total, attack, wait]             1m1s, 1m0s, 674.912ms
Latencies     [min, mean, 50, 90, 95, 99, max]  1.785ms, 492.877ms, 349.044ms, 863.455ms, 1.016s, 1.385s, 2.039s
Bytes In      [total, mean]                     1461793, 48.73
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           25.34%
Status Codes  [code:count]                      200:7603  422:22397
Error Set:
422 Unprocessable Entity
```

### 1:1:1000

baseline:

```
echo "GET http://downstream:3333" | vegeta attack -rate=1000  -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         60000, 1000.02, 818.71
Duration      [total, attack, wait]             1m1s, 59.999s, 551.17ms
Latencies     [min, mean, 50, 90, 95, 99, max]  966.618??s, 300.84ms, 300.635ms, 541.499ms, 571.344ms, 594.501ms, 630.107ms
Bytes In      [total, mean]                     3457793, 57.63
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           82.62%
Status Codes  [code:count]                      200:49573  400:1762  408:1744  500:3483  502:1703  503:1735
Error Set:
500 Internal Server Error
408 Request Timeout
502 Bad Gateway
400 Bad Request
503 Service Unavailable
```

py:
```
echo "GET http://py:9002/py?timeout=333" | vegeta attack -rate=1000 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         58603, 814.39, 0.73
Duration      [total, attack, wait]             1m41s, 1m12s, 28.727s
Latencies     [min, mean, 50, 90, 95, 99, max]  27.512ms, 21.223s, 30.001s, 33.207s, 35.535s, 40.275s, 48.695s
Bytes In      [total, mean]                     7036, 0.12
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           0.13%
Status Codes  [code:count]                      0:58448  200:74  422:79  500:2
Error Set:
422 UNPROCESSABLE ENTITY
...Number of `read: connection reset by peer`...
```

py-ev:
```
echo "GET http://py-ev:9002/py-ev?timeout=333" | vegeta attack -rate=1000 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         56421, 938.12, 0.28
Duration      [total, attack, wait]             1m26s, 1m0s, 25.811s
Latencies     [min, mean, 50, 90, 95, 99, max]  7.467ms, 20.531s, 26.059s, 33.53s, 35.303s, 37.215s, 40.036s
Bytes In      [total, mean]                     106912, 1.89
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           0.04%
Status Codes  [code:count]                      0:49813  200:24  422:6584
Error Set:
422
...Number of `read: connection reset by peer`...
```

ts:
```
echo "GET http://ts:2020/ts?timeout=333" | vegeta attack -rate=1000 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         59250, 873.55, 0.91
Duration      [total, attack, wait]             1m30s, 1m8s, 21.888s
Latencies     [min, mean, 50, 90, 95, 99, max]  6.767ms, 18.797s, 17.743s, 33.219s, 34.017s, 35.365s, 40.717s
Bytes In      [total, mean]                     331266, 5.59
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           0.14%
Status Codes  [code:count]                      0:51589  200:82  422:7579
Error Set:
422 Unprocessable Entity
...Number of `Get "http://ts:2020/ts?timeout=333": dial tcp 0.0.0.0:0->10.196.37.34:2020: bind: address already in use`
...Number of `Get "http://ts:2020/ts?timeout=333": context deadline exceeded (Client.Timeout exceeded while awaiting headers)`
```

go:
```
echo "GET http://go:8090/go?timeout=333" | vegeta attack -rate=1000 -duration=1m -workers=100 | vegeta report
Requests      [total, rate, throughput]         59997, 999.92, 120.45
Duration      [total, attack, wait]             1m1s, 1m0s, 1.219s
Latencies     [min, mean, 50, 90, 95, 99, max]  2.29ms, 698.272ms, 668.993ms, 1.176s, 1.347s, 1.688s, 2.522s
Bytes In      [total, mean]                     2746205, 45.77
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           12.29%
Status Codes  [code:count]                      200:7374  422:52623
Error Set:
422 Unprocessable Entity
```

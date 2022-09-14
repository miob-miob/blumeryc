# Load testing

## All in cluster

* downstream service deployed to cluster
* each implemantion deployed to cluster
* load testing tooling deployed and tests executed inside the cluster

## Results

1 replica

```
echo "GET http://ts:2020/ts?timeout=1000" | vegeta attack -rate=100 -duration=10s | vegeta report
Requests      [total, rate, throughput]         1000, 100.10, 90.83
Duration      [total, attack, wait]             10.636s, 9.99s, 646.027ms
Latencies     [min, mean, 50, 90, 95, 99, max]  7.114ms, 449.589ms, 397.071ms, 871.749ms, 1.046s, 1.285s, 1.479s
Bytes In      [total, mean]                     64876, 64.88
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           96.60%
Status Codes  [code:count]                      200:966  500:34
Error Set:
500 Internal Server Error
```

```
echo "GET http://py-ev:9002/py-ev?timeout=1000" | vegeta attack -rate=100 -duration=10s | vegeta report
Requests      [total, rate, throughput]         998, 99.88, 92.38
Duration      [total, attack, wait]             10.749s, 9.992s, 756.95ms
Latencies     [min, mean, 50, 90, 95, 99, max]  6.762ms, 375.828ms, 362.411ms, 664.684ms, 763.406ms, 945.641ms, 1.149s
Bytes In      [total, mean]                     65446, 65.58
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           99.50%
Status Codes  [code:count]                      200:993  500:5
```

```
echo "GET http://go:8090/go?timeout=1000" | vegeta attack -rate=100 -duration=10s | vegeta report
```

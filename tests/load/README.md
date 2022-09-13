# Load testing

## Stack

Vegeta load test tool https://github.com/tsenart/vegeta

```
echo "GET http://localhost:8080" | vegeta attack -rate=128 -duration=10s | vegeta report
```

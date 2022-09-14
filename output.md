
# Output
 

## how the selected technologies supported parallelism (and concurrency),

- GO
    - d
- TS
	- event loop 
	- single thread
- PY
	- depends, if we use code based on [asyncio](https://docs.python.org/3/library/asyncio.html) (asgi) it will be similar like node js
	- if not we will have classic thread /request (apache + mod_wsgi) 
	
conclusion: 
    it really depends on what does __high__ performance really means. If we can afford 
    node, development may be cheaper. Weakest point of application also matters  - in our case it appears to be downstream service 
## How will the different approaches in parallelism (and concurrency) affect the solution's scalability?


## What tools will you use for autoscaling? 


------------------------------------------ 

## Please consider Kubernetes (GKE) 

- Considered


## Suggest the best way for resource allocation and how we may solve the autoscaling


## Will you prefer horizontal or vertical autoscaling?

## What metrics will you use as inputs for autoscaling?

## simple API performance tests which

- DONE
- TODO: add output metrics
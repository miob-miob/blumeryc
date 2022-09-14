
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
	
conclusion: 

it really depends on what does __high__ performance really means. If we can afford 
node, development cost may be cheaper. Weakest point of application also matters  - in our case it appears to be downstream service 

In this case the pick of the technology does not mean much,
because the bottleneck of this app was in the third part party downstream service.

If we'll have more CPU heavy loaded algorithms in our codebase, 
GO will have much better performance in comparision with interpreted languages.

**fun fact:**
For this kind of assignement we picked python with event loop to not to waste resources with creating new threds which will wait till thre reponse from downstream service will come.

## How will the different approaches in parallelism (and concurrency) affect the solution's scalability?

(viz @jurasek charts)

------------------------------------------ 

## What tools will you use for autoscaling? 

## Please consider Kubernetes (GKE) 

- Considered


## Suggest the best way for resource allocation and how we may solve the autoscaling


## Will you prefer horizontal or vertical autoscaling?

## What metrics will you use as inputs for autoscaling?

## simple API performance tests which

- DONE
- TODO: add output metrics
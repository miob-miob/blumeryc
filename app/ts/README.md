# TS implementation

## Run local docker instance

```sh
npm run dev

docker build . -t tsahoj

docker run -p 2020:2020 \
	-e DOWNSTREAM_URL=https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp \
	tsahoj

curl http://localhost:2020/ts\?timeout\=800


```

## docker hub

https://hub.docker.com/repository/docker/nevimjakymamid/tsahoj

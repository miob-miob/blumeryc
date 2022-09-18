# TS implementation

## Local development

```sh
npm run dev
```

## Docker prod

```
docker build . -t nevimjakymamid/tsahoj

docker run -p 2020:2020 \
	-e DOWNSTREAM_URL=https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp \
	nevimjakymamid/tsahoj

curl http://localhost:2020/ts\?timeout\=800


```

## docker hub

https://hub.docker.com/repository/docker/nevimjakymamid/tsahoj

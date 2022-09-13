## Downstream service

PROD API URL:

https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp?timeout=20

## lambda configuration

lambda timeout is 60sec in max

```sh
npm run dev

docker build . -t tsahoj

docker run -p 2020:2020 \
	-e DOWNSTREAM_URL=https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp \
	tsahoj

curl http://localhost:2020/ts\?timeout\=800

```

## GO implementation

```sh
DOWNSTREAM_URL=https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp go run index.go


docker build . -t nevimjakymamid/go-dominik

docker run -p 8090:8090 \
	-e DOWNSTREAM_URL=https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp \
  nevimjakymamid/go-dominik

curl http://localhost:8090/go\?timeout\=800

```
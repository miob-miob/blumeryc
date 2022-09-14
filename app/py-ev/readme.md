###Local development

```
virtualenv venv -p python3.8.9
cd venv
source bin/activate
cd .. 
pip install -r requirements.txt

export DOWNSTREAM_URL="https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp"
python server.py

curl localhost:9002/py-ev\?timeout\=1000

```

### docker image

[docker hub link](https://hub.docker.com/repository/docker/oborim01/kukurice)
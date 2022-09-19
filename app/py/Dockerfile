FROM python:3.8-alpine

WORKDIR app
COPY requirements.txt requirements.txt

RUN pip install -r requirements.txt \
    pip install pyuwsgi



COPY server.py server.py


EXPOSE 9002
ENV RUN_PROD 'prod'
#ENTRYPOINT ["hypercorn", "-b", "0.0.0.0:9002", "server:app"]
ENTRYPOINT uwsgi --http 0.0.0.0:9002 --master -p 1 -w server:app

# docker run -it --rm -p 9002:9002 -e DOWNSTREAM_URL=https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp oborim01/non_kokurice

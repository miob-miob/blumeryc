FROM python:3.8-alpine

WORKDIR app
COPY requirements.txt requirements.txt

RUN pip install -r requirements.txt

COPY server.py server.py

EXPOSE 9002
ENV RUN_PROD 'prod'
ENTRYPOINT ["hypercorn", "-b", "0.0.0.0:9002", "server:app"]

# docker run -it --rm -p 9002:9002 -e DOWNSTREAM_URL="http://192.168.0.129:3333" miob:kukurice

# pip install -r /path/to/requirements.txt
# hypercorn -b 0.0.0.0:9002 server:app

FROM python:3.8-alpine

WORKDIR APP
COPY requirements.txt requirements.txt
COPY server.py server.py


RUN pip install -r requirements.txt


EXPOSE 9002
ENV RUN_PROD 'prod'
ENTRYPOINT ["hypercorn", "-b", "0.0.0.0:9002", "server:app"]
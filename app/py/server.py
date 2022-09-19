from threading import Thread
import requests
from time import sleep
from requests.exceptions import ConnectTimeout, ReadTimeout
from queue import Queue
from flask import Flask, Response, request
import os

DOWNSTREAM_URL = os.environ.get('DOWNSTREAM_URL')
PORT = 9002

if DOWNSTREAM_URL is None:
    raise ValueError("'DOWNSTREAM_URL environment variable is missing!")

DOOMS_SERVER_ON_LAMBDA = 'https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp'
FIRST_TIMEOUT = 0.3

PRODUCTION = True if os.environ.get('RUN_ENV') == 'prod' else False


class NotCorrectResponseError(ValueError):
    pass


class HardTimeOut(TimeoutError):
    pass


class FirstTimeout(TimeoutError):
    pass


def validate_response(data: dict):
    errors = []
    if data.get('requestId') is None:
        errors.append('missing requestId on downstream server response')
    if data.get('timeout') is None:
        errors.append('missing timeout on downstream server response')
    return errors


def perform_request(url, hard_timeout):
    try:
        # this timeout forces thread to die within some meaningful limit
        res = requests.get(url, timeout=hard_timeout)
    except (ReadTimeout, ConnectTimeout):
        raise NotCorrectResponseError('Timeout occurred!')
    if res.status_code != 200:
        raise NotCorrectResponseError('Response status != 200!!!')
    try:
        json_data = res.json()
    except Exception:
        raise NotCorrectResponseError('JSON parsing died')
    errors = validate_response(json_data)
    if len(errors) > 0:
        raise NotCorrectResponseError(' '.join(errors))
    return json_data


def thread_fn(fn, args, queue: Queue):
    try:
        result = fn(*args)
        queue.put_nowait(result)
    except NotCorrectResponseError as e:
        queue.put_nowait(e)


def first_timeout_fn(queue: Queue):
    sleep(FIRST_TIMEOUT)
    queue.put_nowait(FirstTimeout())


def hard_timeout_fn(queue: Queue, hard_timeout):
    sleep(hard_timeout)
    queue.put_nowait(HardTimeOut())


app = Flask(__name__)


@app.route("/py")
def das_endpoint():
    hard_timeout = request.args.get('timeout')
    if hard_timeout is None:
        return Response("Missing timeout", 400)
    try:
        hard_timeout = int(hard_timeout) / 1000
    except ValueError:
        return Response("timeout must be parsable as int", 400)
    if hard_timeout < FIRST_TIMEOUT:
        return Response(f'timeout cannot be smaller then {FIRST_TIMEOUT}', 400)

    main_q = Queue()

    # hard timeout
    t0 = Thread(target=hard_timeout_fn, args=[main_q, hard_timeout], daemon=True)
    t0.start()

    # first request
    t1 = Thread(target=thread_fn, args=[perform_request, [DOWNSTREAM_URL,hard_timeout], main_q], daemon=True)
    t1.start()

    # first timeout
    t2 = Thread(target=first_timeout_fn, args=[main_q], daemon=True)
    t2.start()

    res = main_q.get()

    max_fails = 3
    fails = 0
    # if first failed prior timeout
    failed_prior_timeout = isinstance(res, NotCorrectResponseError)

    first_timeout_occurred = isinstance(res, FirstTimeout)
    # we will need extra requests
    if first_timeout_occurred or failed_prior_timeout:
        print('we will need more requests')
        if failed_prior_timeout:
            print('quick fail!')
            fails += 1
        t3 = Thread(target=thread_fn, args=[perform_request, [DOWNSTREAM_URL,hard_timeout], main_q], daemon=True)
        t4 = Thread(target=thread_fn, args=[perform_request, [DOWNSTREAM_URL,hard_timeout], main_q], daemon=True)
        t3.start()
        t4.start()
    # first was successful within time limit - we are done
    else:
        print('First one was success!')
        return res

    while fails < max_fails:
        result = main_q.get()
        if isinstance(result, HardTimeOut):
            print('Hard timeout reached')
            return Response('Hard timeout reached', 422)
        # ignore first timer (in case of quick fail)
        elif isinstance(result, FirstTimeout):
            pass
        elif isinstance(result, NotCorrectResponseError):
            print('One of requests failed!')
            fails += 1
        else:
            print('One of request was successful')
            return result

    print("None of request was successful!!!")
    return Response("None of request was successful!!!", 422)


if __name__ == '__main__':

    if PRODUCTION:
        print("you should use real WSGI server")

    else:
        print('starting dev server')
        app.run(host='0.0.0.0', port=PORT, use_reloader=True, debug=True)

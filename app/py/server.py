from threading import Thread
import requests
from time import sleep
from requests.exceptions import ConnectTimeout, ReadTimeout
from queue import Queue

down_stream_url = 'https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp'
hard_timeout = 1
first_timeout = 0.3


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


def perform_request(url=down_stream_url, timeout=hard_timeout):
    try:
        res = requests.get(url, timeout=timeout)
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
    sleep(first_timeout)
    queue.put_nowait(FirstTimeout())


def hard_timeout_fn(queue: Queue):
    sleep(hard_timeout)
    queue.put_nowait(HardTimeOut())


def main():
    main_q = Queue()

    # hard timeout
    t0 = Thread(target=hard_timeout_fn, args=[main_q],daemon=True)
    t0.start()

    # first request
    t1 = Thread(target=thread_fn, args=[perform_request, [], main_q],daemon=True)
    t1.start()

    # first timeout
    t2 = Thread(target=first_timeout_fn, args=[main_q],daemon=True)
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
            fails +=1
        t3 = Thread(target=thread_fn, args=[perform_request, [], main_q], daemon=True)
        t4 = Thread(target=thread_fn, args=[perform_request, [], main_q], daemon=True)
        t3.start()
        t4.start()
    # first was successful within time limit - we are done
    else:
        print('Successful and quick first one!', res)
        return



    while fails < max_fails:
        result = main_q.get()
        if isinstance(result, HardTimeOut):
            print('Hard timeout reached')
            return
        # ignore first timer (in case of quick fail)
        elif isinstance(result, FirstTimeout):
            pass
        elif isinstance(result,NotCorrectResponseError):
            print('some fail!')
            fails += 1
        else:
            print('Success',result)
            return

    print("None of request was successful!!!")


if __name__ == '__main__':
    main()

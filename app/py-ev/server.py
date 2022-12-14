from quart import Quart, request, Response
import asyncio
import os
import aiohttp

DOWNSTREAM_URL = os.environ.get('DOWNSTREAM_URL')
PORT = 9002
TIME_FOR_FIRST = 0.3

if DOWNSTREAM_URL is None:
    raise ValueError("'DOWNSTREAM_URL environment variable is missing!")

app = Quart(__name__)



def validate_response(data: dict):
    errors = []
    if data.get('requestId') is None:
        errors.append('missing requestId on downstream server response')
    if data.get('timeout') is None:
        errors.append('missing timeout on downstream server response')
    return errors


class NotCorrectResponseError(ValueError):
    pass


async def perform_request(http_client_session: aiohttp.ClientSession, url=DOWNSTREAM_URL):
    r = await http_client_session.get(url)
    if r.status == 200:
        try:
            json_data = await r.json()
        # to broad, i know, but simple ;)
        except Exception:
            return NotCorrectResponseError("NOt able to parse JSON from response")
        errors = validate_response(json_data)
        if len(errors) != 0:
            raise NotCorrectResponseError(' ,'.join(errors))
        return json_data

    else:
        return NotCorrectResponseError("Down stream failed")


@app.route('/py-ev')
async def das_endpoint():
    timeout = request.args.get('timeout')
    if timeout == None:
        return Response("Missing timeout", 400)

    try:
        timeout = int(timeout) / 1000
    except ValueError:
        return Response("timeout must be parsable as int", 400)
    if timeout < TIME_FOR_FIRST:
        return Response(f'timeout cannot be smaller then {TIME_FOR_FIRST}', 400)
    async with aiohttp.ClientSession() as session:
        first_request = asyncio.create_task(perform_request(session))
        done, pending = await asyncio.wait({first_request}, timeout=TIME_FOR_FIRST, return_when=asyncio.ALL_COMPLETED)
        first_processed = False

        if first_request in done:
            first_processed = True
            data = await first_request
            if not isinstance(data, NotCorrectResponseError):
                return data

        second_request = asyncio.create_task(perform_request(session))
        third_request = asyncio.create_task(perform_request(session))
        requests = [second_request, third_request]

        if not first_processed:
            requests.append(first_request)
        try:
            for coro in asyncio.as_completed(requests, timeout=timeout - TIME_FOR_FIRST):
                next_done = await coro
                if not isinstance(next_done, NotCorrectResponseError):
                    for task in requests:
                        task.cancel()
                    return next_done

        except asyncio.TimeoutError:
            for task in requests:
                task.cancel()

            return Response('Timeout reached!', 422)

        return Response("Down stream service is not working properly", 422)


if __name__ == '__main__':
    production = True if os.environ.get('RUN_ENV') == 'prod' else False


    if production:
        # hypercorn -b 0.0.0.0:9002 server:app
        print(f'you should run `hypercorn -b 0.0.0.0:{PORT} server:app`')

    else:
        print('starting dev server')
        app.run(host='0.0.0.0', port=PORT, use_reloader=True, debug=False)

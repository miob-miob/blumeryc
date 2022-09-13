/* eslint-disable @typescript-eslint/no-var-requires */
require('dotenv').config()
import 'express-async-errors'
import { getNumberFromEnvParser, validateConfig } from 'typed-env-parser'
import cors from 'cors'
import express from 'express'
import fetch from 'node-fetch'

export const appEnvs = validateConfig({
  PORT: getNumberFromEnvParser('PORT'),
})

export const appConfig = {
  downstreamServiceURL:
    // eslint-disable-next-line max-len
    'https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp',
  DOWNSTREAM_SERVICE_TIMEOUT_MS: 300,
} as const

type DownstreamServiceData = {
  requestId: string
  timeout: number
}

const delay = (ms: number) => new Promise(res => setTimeout(res, ms))

const fetchWithTimeout = async (resource: string, options: { timeout: number }) => {
  const controller = new AbortController()
  const id = setTimeout(() => controller.abort(), options.timeout)
  const response = await fetch(resource, {
    ...options,
    // @ts-expect-error ???
    signal: controller.signal,
  })
  clearTimeout(id)
  return response
}

const services = {
  // in the JavaScript fetch API you're not able to close HTTP connections
  // it throws errors even when is ok but body is missing
  // TODO: add timeout for fetch requests
  getDownstreamData: async (options: { timeout: number }) => {
    const response = await fetchWithTimeout(appConfig.downstreamServiceURL, options)

    if (!response.ok) throw new Error('invalid HTTP network call ' + response)
    const data = (await response.json()) as DownstreamServiceData

    const isValid =
      response.status === 200 &&
      typeof data.requestId === 'string' &&
      typeof data.timeout === 'number'

    if (!isValid) throw new Error('Invalid shape for success response')

    return data
  },
}

const app = express()
app.use(express.urlencoded({ extended: true }))
app.use(express.json())
app.use(cors())

// assignment:
//   "The endpoint always returns a response within the given timeout."
// response:
//   "There is a little overhead of running javascript which care about promise resolving"
app.get('/ts', async (req, res) => {
  console.info('calling ts endpoint')
  // stop working care to not to the redundant HTTP requests
  let stopWorking = false

  const qTimeout = typeof req.query.timeout === 'string' ? parseInt(req.query.timeout, 10) : NaN

  if (isNaN(qTimeout)) {
    res.status(400).send('invalid timeout parameter')
    return
  }
  if (qTimeout <= appConfig.DOWNSTREAM_SERVICE_TIMEOUT_MS) {
    res.status(400).send(`timeout parameter has to be > ${appConfig.DOWNSTREAM_SERVICE_TIMEOUT_MS}`)
    return
  }

  try {
    const data = await Promise.any([
      services.getDownstreamData({ timeout: appConfig.DOWNSTREAM_SERVICE_TIMEOUT_MS }),

      // fetch 2. API call
      (async () => {
        await delay(appConfig.DOWNSTREAM_SERVICE_TIMEOUT_MS)
        if (stopWorking) return
        return services.getDownstreamData({
          timeout: qTimeout - appConfig.DOWNSTREAM_SERVICE_TIMEOUT_MS,
        })
      })(),

      // fetch 3. API call
      (async () => {
        await delay(appConfig.DOWNSTREAM_SERVICE_TIMEOUT_MS)
        if (stopWorking) return
        return services.getDownstreamData({
          timeout: qTimeout - appConfig.DOWNSTREAM_SERVICE_TIMEOUT_MS,
        })
      })(),
    ])

    stopWorking = true

    res.json(data)
  } catch (err) {
    console.error('err ', Math.random())
    // All promises were rejected => should never happen...
    res.status(408).send(`downstream services is not working properly`)
  }
})

app.get('*', (_req, res) => {
  res.status(404).send(`
    <h1>route not found</h1>
    <h2>check out the rýč</h2>
    <img src="https://github.com/miob-miob/blumeryc/raw/master/logo.png"></img>`)
})

app.listen(appEnvs.PORT, () => {
  console.info(`
--------- server is ready now ---------
URL: http://localhost:${appEnvs.PORT}/ts
---------------------------------------
  `)
})

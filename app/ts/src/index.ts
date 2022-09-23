/* eslint-disable @typescript-eslint/no-var-requires */
import 'express-async-errors'
import { getStringFromEnvParser, validateConfig } from 'typed-env-parser'
import cors from 'cors'
import express from 'express'
import fetch from 'node-fetch'
import process from 'process'

// keep killing process on CMD+C if the app is running in the docker
process.on('SIGINT', () => {
  console.info('💀💀💀💀️️')
  process.exit()
})

export const appEnvs = validateConfig({
  downstreamServiceURL: getStringFromEnvParser('DOWNSTREAM_URL'),
})

export const appConfig = {
  DOWNSTREAM_SERVICE_TIMEOUT_MS: 300,
  PORT: 2020,
} as const

type DownstreamServiceData = {
  requestId: string
  timeout: number
}

const delay = (ms: number) => new Promise(res => setTimeout(res, ms))

const services = {
  getDownstreamData: async () => {
    const response = await fetch(appEnvs.downstreamServiceURL)

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
    const initReq = services.getDownstreamData()

    const throwAfterGlobalTimeoutPromise = (async () => {
      await delay(qTimeout)
      // we have to handle that we will not throw error for sent request
      // because it will not be caught by express handler try catch and error is propagated to the unhandledRejections
      if (!res.writableEnded) throw new Error('GLOBAL_TIMEOUT_EXCEEDED')
    })()

    try {
      const initReqOKResponse = await Promise.race([
        initReq,
        (async () => {
          await delay(appConfig.DOWNSTREAM_SERVICE_TIMEOUT_MS)
          throw new Error('TIMEOUT_EXCEEDED')
        })(),
      ])

      res.json(initReqOKResponse)
      return
    } catch (err) {
      // continue fetching...
    }

    const data = await Promise.race([
      throwAfterGlobalTimeoutPromise,

      Promise.any([
        // keep fetching 1st API call
        initReq,

        // fetch 2nd API call
        services.getDownstreamData(),

        // fetch 3rd API call
        services.getDownstreamData(),
      ]),
    ])

    res.json(data!)
  } catch (err) {
    res.status(422).send(`downstream services is not working properly`)
  }
})

app.get('*', (_req, res) => {
  res.status(404).send(`
    <h1>route not found</h1>
    <h2>check out the rýč</h2>
    <img src="https://github.com/miob-miob/blumeryc/raw/master/logo.png"></img>`)
})

app.listen(appConfig.PORT, () => {
  console.info(`
--------- server is ready now ---------
URL: http://localhost:${appConfig.PORT}/ts
---------------------------------------
  `)
})

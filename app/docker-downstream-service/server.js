const process = require('process')
const express = require("express");
const { v4: uuidv4 } = require('uuid');

const PORT = process.env.PORT || 3333;
const FAILURE_RATE = process.env.FAILURE_RATE || 0.2; // 20%

// keep killing process on CMD+C if the app is running in the docker
process.on('SIGINT', () => {
  console.info('💀💀💀💀️️')
  process.exit(0)
})

const app = express();

app.use(express.json());

app.get('/', (req, res) => {
   const timeout = Math.round(Math.random() * 600)

   const failures = [
      {responseCode:200},
      {responseCode:400},
      {responseCode:408},
      {responseCode:500},
      {responseCode:500, responseData: {requestId: uuidv4(), timeout: 333}},
      {responseCode:502},
      {responseCode:503},
   ]
   const failure = Math.random() < FAILURE_RATE ? failures[Math.round(Math.random() * failures.length)] : null
   return setTimeout(
      _=> {
         if (failure) {
            if (failure.responseData) {
               res
                  .status(failure.responseCode)
                  .json(failure.responseData);
            } else {
               res.sendStatus(failure.responseCode)
            }
         } else {
            res
               .status(200)
               .json({
                  requestId: uuidv4(),
                  timeout: timeout
               });
         }
      },
      timeout
   )
   
});

app.listen(PORT, () => console.log(`Assignment testing server listening on port ${PORT}!`));
const serverless = require('serverless-http')
const express = require('express')

const PORT = process.env.PORT || 3333;
const FAILURE_RATE = process.env.FAILURE_RATE || 0.2; // 20%

// const express = require("express");
const { v4: uuidv4 } = require('uuid');

const app = express();

app.use(express.json());

app.get('/default/blumeryc-downstream-service-dominik-tilp', (req, res) => {

   // missing implementation for req param in the previous implementations
   const queryTimeout = parseInt(req.query.timeout, 10)

   const timeout = !isNaN(queryTimeout)
      ? queryTimeout
      : 600

   const reqTimeout = Math.round(Math.random() * timeout)

   if (timeout > 59_000) {
      res.send(400).send('invalid timeout parameter, has to be <= 59sec')
      return
   }

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
      reqTimeout
   )
   
});

// app.listen(PORT, () => console.log(`Assignment testing server listening on port ${PORT}!`));

module.exports.handler = async (...args) => {
   return serverless(app)(...args)
}
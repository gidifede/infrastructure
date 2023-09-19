const express = require('express')
const serverless = require('serverless-http')
const swaggerUI = require('swagger-ui-express')
const swaggerDocument = require('/opt/swagger.json');


const app = express()

module.exports.handler = async (event, context) => {

    const API_PREFIX = process.env.API_PREFIX

    app.use(API_PREFIX + '/api-docs', swaggerUI.serve, swaggerUI.setup(swaggerDocument))
    const handler = serverless(app);

    return await handler(event, context)
}

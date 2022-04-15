'use strict';

// disable for entire file: we dynamically require modules based on `appConfig`
/* eslint-disable global-require */

//
// Web Application (served from the web server)
//

// Library imports
const express = require('express');
const rateLimit = require("express-rate-limit");
const busboy = require('connect-busboy');
const path = require('path');
const mkdirp = require('mkdirp');
const morgan = require('morgan');
const makeRfs = require('rotating-file-stream');
const winston = require('winston');

// App imports
const paths = require('./conf/paths');
const appConfig = require(paths.src.appConfig);
const libUtil = require(paths.src.util);
const bodyParser = require(paths.src.bodyParser);

//
// Create and configure the app server
//
const app = express();

app.use(
    rateLimit({
        windowMs: 1000, // 1 second duration
        max: 20,
        message: "You exceeded 20 requests per second; to increase the allowed requests/sec, modify rateLimit in app.js",
        headers: true
    })
);

winston.verbose('Mounting route to handle server error state');
app.use((req, res, next) => {
    winston.debug('Checking for error state', {state: appConfig.serverErrorState});
    if (appConfig.serverErrorState) {
        winston.error('Server is in error state; serving up error html', {
            serverErrorState: appConfig.serverErrorState
        });
        res.sendFile(paths.resources.initErrorHtml);
    } else {
        next();
    }
});

winston.info('--- app server initializing ---');

app.use(busboy({
    highWaterMark: 2 * 1024 * 1024, // Set 2MiB buffer
})); // Insert the busboy middle-ware

// configure request logging
(function() {
    // create if doesn't exist
    mkdirp.sync(appConfig.logPath);

    // rotating file stream filename generator
    const generator = libUtil.makeRotatingFilenameGenerator(
        path.join(appConfig.logPath, 'access.log'));

    // rotate every 5 megabytes, append if file exists
    const stream = makeRfs(generator, {
        size: '5M'
    });

    app.use(morgan(
        ':date[iso] :remote-addr :remote-user ":method :url HTTP/:http-version" ' +
        ':status :res[content-length] - :response-time ms - ":user-agent"', {
            stream
        }));
})();

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: false
}));

if (appConfig.clientPath) {
    winston.info('client static path: %s', appConfig.clientPath);

    // assume this path provides an index.html and provide it at the root
    app.use('/', express.static(appConfig.clientPath));
}

// local API routes are used, and they load the mock rest API services
winston.info('using API mock REST services');

const restApiRoutes = require(paths.src.routes.common);
const dockerApiRoutes = require(paths.src.routes.docker);
const vsphereApiRoutes = require(paths.src.routes.vsphere);
const awsApiRoutes = require(paths.src.routes.aws);
const azureApiRoutes = require(paths.src.routes.azure);

const ENDPOINT = appConfig.apiEndpoint;

app.use(`${ENDPOINT}`, restApiRoutes);
app.use(`${ENDPOINT}/providers/docker`, dockerApiRoutes);
app.use(`${ENDPOINT}/providers/vsphere`, vsphereApiRoutes);
app.use(`${ENDPOINT}/providers/aws`, awsApiRoutes);
app.use(`${ENDPOINT}/providers/azure`, azureApiRoutes);


// Catch all other routes and return the index file.
// Let the client application handle the http status code
// 404 situation in case of non-existent route.
app.get('*', (req, res) => {
    res.sendFile(appConfig.clientPath + '/index.html');
});

//// error handlers

// Since middleware are executed in the order they are added (via `app.use(...)`),
//  and since all valid routes MUST have been added by this point in app initialization,
//  if this middleware executes, it essentially constitutes a 404 because no other
//  middleware handled the request: generate a 404 and forward to error handlers.
app.use((req, res, next) => {
    const err = new Error('Not Found');
    err.status = 404;
    next(err);
});

// development error handler: will print stacktrace
app.use((err, req, res, next) => {
    const status = err.status || 500;
    res.status(status);
    res.send({
        error: err,
        message: err.message
    });
});

//
// Module definition
//
module.exports = app;

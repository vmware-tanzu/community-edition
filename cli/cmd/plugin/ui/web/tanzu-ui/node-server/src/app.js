'use strict';

// disable for entire file: we dynamically require modules based on `appConfig`
/* eslint-disable global-require */

//
// Web Application (served from the web server)
//

const express = require('express');
const rateLimit = require("express-rate-limit");

const busboy = require('connect-busboy');
const path = require('path');
const paths = require('./conf/paths');
const mkdirp = require('mkdirp');
const morgan = require('morgan');
const makeRfs = require('rotating-file-stream');
const appConfig = require(paths.src.appConfig);
const libUtil = require(paths.src.util);
const winston = require('winston');

const bodyParser = require(paths.src.bodyParser);

//
// Create and configure the app server
//
let app = express();

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
    let generator = libUtil.makeRotatingFilenameGenerator(
        path.join(appConfig.logPath, 'access.log'));

    // rotate every 5 megabytes, append if file exists
    let stream = makeRfs(generator, {
        size: '5M'
    });

    app.use(morgan(
        ':date[iso] :remote-addr :remote-user ":method :url HTTP/:http-version" ' +
        ':status :res[content-length] - :response-time ms - ":user-agent"', {
            stream
        }));
})();

let sessionConfig = {
    secret: 'tceNodeSession',
    rolling: true,
    resave: true,
    saveUninitialized: true,
    name: 'session',
    cookie: {
        secure: true,
        httpOnly: true
    }
};

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
let restApiRoutes = require(paths.src.restApiRoutes);
app.use(restApiRoutes);


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
    let err = new Error('Not Found');
    err.status = 404;
    next(err);
});

// development error handler: will print stacktrace
app.use((err, req, res, next) => {
    let status = err.status || 500;
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

'use strict';

//
// Web Server
//
const _ = require('lodash');
const fs = require('fs');
const http = require('http');
const https = require('https');
const paths = require('./conf/paths');
const app = require(paths.src.app);
const appConfig = require(paths.src.appConfig);
const winston = require('winston');

const WebSocket = require('ws').Server;

// {Number|String} port number or named pipe
let port;
// {Node|http}
let server;

// Normalize a port into a number, string, or false.
// @returns {Number|String|Boolean} Normalized port as a string (named pipe), number
//   (port), or `false` (invalid|unknown).
let normalizePort = function(val) {
    let normPort = parseInt(val, 10);

    if (isNaN(normPort)) {
        // named pipe
        return val;
    }

    if (normPort >= 0) {
        // port number
        return normPort;
    }

    return false;
};

// Event listener for HTTP server 'error' event.
// @param {Node|http.error} error
let onError = function(error) {
    let bind = _.isString(port) ? ('pipe ' + port) : ('port ' + port);

    if (error.syscall !== 'listen') {
        throw error;
    }

    // handle specific listen errors with friendly messages
    switch (error.code) {
        case 'EACCES':
            winston.error(bind + ' requires elevated privileges');
            process.exit(1);
            break;

        case 'EADDRINUSE':
            winston.error(bind + ' is already in use');
            process.exit(1);
            break;

        default:
            throw error;
    }
};

// Event listener for HTTP server 'listening' event.
let onListening = function() {
    let addr = server.address();
    let bind = _.isString(addr) ? ('pipe ' + addr) : ('port ' + addr.port);
    winston.info('listening on ' + bind);
};

// get port from environment and store in Express
port = normalizePort(appConfig.port);
app.set('port', port);

// create HTTP server
winston.warn('Starting UI server in insecure (http) mode');
server = http.createServer(app);

const logData = [
    {
        type: 'status',
        data: {
            message: 'start deploying',
            status: 'running',
            currentPhase: 'Configure prerequisite',
            totalPhases: ['Configure prerequisite',
                'Validate configuration',
                'Generate cluster configuration',
                'Setup bootstrap cluster',
                'Install providers on bootstrap cluster',
                'Create management cluster',
                'Install providers on management cluster',
                'Move cluster-api objects from bootstrap cluster to management cluster'],
            successfulPhases: []
        }
    },
    {
        type: 'log',
        data: {
            currentPhase: 'Deploy Cluster phase msg 1',
            message: 'I0227 16:23:18.924759 ms1',
            logType: 'ERROR'
        }
    },
    {
        type: 'log',
        data: {
            currentPhase: 'Deploy Cluster phase msg 2',
            message: 'I0227 16:23:18.924759 msg2',
            logType: 'INFO'
        }
    },
    {
        type: 'log',
        data: {
            currentPhase: 'Deploy Cluster phase msg 3',
            message: 'I0227 16:23:18.924759 msg3',
            logType: 'INFO'
        }
    },
    {
        type: 'status',
        data: {
            message: 'start deploying',
            status: 'running',
            currentPhase: 'Validate configuration',
            totalPhases: ['Configure prerequisite',
                'Validate configuration',
                'Generate cluster configuration',
                'Setup bootstrap cluster',
                'Install providers on bootstrap cluster',
                'Create management cluster',
                'Install providers on management cluster',
                'Move cluster-api objects from bootstrap cluster to management cluster'],
            successfulPhases: []
        }
    },
    {
        type: 'status',
        data: {
            message: 'start deploying',
            status: 'running',
            currentPhase: 'Generate cluster configuration',
            totalPhases: ['Configure prerequisite',
                'Validate configuration',
                'Generate cluster configuration',
                'Setup bootstrap cluster',
                'Install providers on bootstrap cluster',
                'Create management cluster',
                'Install providers on management cluster',
                'Move cluster-api objects from bootstrap cluster to management cluster'],
            successfulPhases: []
        }
    },
    {
        type: 'status',
        data: {
            message: 'start deploying',
            status: 'running',
            currentPhase: 'Setup bootstrap cluster',
            totalPhases: ['Configure prerequisite',
                'Validate configuration',
                'Generate cluster configuration',
                'Setup bootstrap cluster',
                'Install providers on bootstrap cluster',
                'Create management cluster',
                'Install providers on management cluster',
                'Move cluster-api objects from bootstrap cluster to management cluster'],
            successfulPhases: ['Move cluster-api objects from bootstrap cluster to management cluster']
        }
    },
    {
        type: 'status',
        data: {
            message: 'start deploying',
            status: 'running',
            currentPhase: 'Install providers on bootstrap cluster',
            totalPhases: ['Configure prerequisite',
                'Validate configuration',
                'Generate cluster configuration',
                'Setup bootstrap cluster',
                'Install providers on bootstrap cluster',
                'Create management cluster',
                'Install providers on management cluster',
                'Move cluster-api objects from bootstrap cluster to management cluster'],
            successfulPhases: ['Move cluster-api objects from bootstrap cluster to management cluster']
        }
    },
    {
        type: 'status',
        data: {
            message: 'start deploying',
            status: 'running',
            currentPhase: 'Create management cluster',
            totalPhases: ['Configure prerequisite',
                'Validate configuration',
                'Generate cluster configuration',
                'Setup bootstrap cluster',
                'Install providers on bootstrap cluster',
                'Create management cluster',
                'Install providers on management cluster',
                'Move cluster-api objects from bootstrap cluster to management cluster'],
            successfulPhases: ['Move cluster-api objects from bootstrap cluster to management cluster']
        }
    },
    {
        type: 'status',
        data: {
            message: 'start deploying',
            status: 'running',
            currentPhase: 'Install providers on management cluster',
            totalPhases: ['Configure prerequisite',
                'Validate configuration',
                'Generate cluster configuration',
                'Setup bootstrap cluster',
                'Install providers on bootstrap cluster',
                'Create management cluster',
                'Install providers on management cluster',
                'Move cluster-api objects from bootstrap cluster to management cluster'],
            successfulPhases: ['Move cluster-api objects from bootstrap cluster to management cluster']
        }
    },
    {
        type: 'status',
        data: {
            message: 'start deploying',
            status: 'running',
            currentPhase: 'Move cluster-api objects from bootstrap cluster to management cluster',
            totalPhases: ['Configure prerequisite',
                'Validate configuration',
                'Generate cluster configuration',
                'Setup bootstrap cluster',
                'Install providers on bootstrap cluster',
                'Create management cluster',
                'Install providers on management cluster',
                'Move cluster-api objects from bootstrap cluster to management cluster'],
            successfulPhases: ['Move cluster-api objects from bootstrap cluster to management cluster']
        }
    },
    {
        type: 'status',
        data: {
            message: 'start deploying',
            status: 'successful',
            totalPhases: ['Configure prerequisite',
                'Validate configuration',
                'Generate cluster configuration',
                'Setup bootstrap cluster',
                'Install providers on bootstrap cluster',
                'Create management cluster',
                'Install providers on management cluster',
                'Move cluster-api objects from bootstrap cluster to management cluster'],
            successfulPhases: ['Move cluster-api objects from bootstrap cluster to management cluster']
        }
    }
];

// set up websocket connection piggy-backed on express router at path '/ws'
const ws = new WebSocket({ server: server, path: '/ws' });

// wire websocket handlers
ws.on('connection', function (ws) {
    // show the connection has been established in the console
    winston.info("WS Connection established:");

    // wire the event handlers
    ws.on('message', function (data) {
        // show the message object in the console
        var message = JSON.parse(data);
        winston.info("WS Message received from client:");
        winston.info(message.operation);

        // send response to received message
        if (message.operation && message.operation === 'logs') {
            let x = 0;

            // mock an interval between sending log data
            setInterval(function() {
                if (x < logData.length) {
                    ws.send(JSON.stringify(logData[x]));
                }
                else return;

                x++;
            }, 1500);

        } else {
            let response = {
                source: "WebAppsNodeJs Application (server)",
                message:"Client Message Received!"
            };

            ws.send(JSON.stringify(response));
        }
    });
});

// listen on provided port, on all network interfaces
server.listen(port);
server.on('error', onError);
server.on('listening', onListening);

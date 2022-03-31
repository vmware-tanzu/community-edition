'use strict';
/* globals __dirname */

// Library imports
const _ = require('lodash');
const mkdirp = require('mkdirp');
const winston = require('winston');

// App imports
const paths = require('./src/conf/paths');
const appConfig = require(paths.src.appConfig);

const help = `
TCE UI Mock API Server

Available options:

  --port <port|pipe>  : Specify an alternate port (number >= ${appConfig.minPort}) or a named pipe
                            (string). Default is ${appConfig.port}.
  --help              : Print this message.
`;

appConfig.clientPath = paths.directories.clientDistDir;

winston.configure({
    transports: [
        new (winston.transports.Console)({
            level: appConfig.logLevel
        })
    ]
});

if (_.includes(process.argv, '--help')) {
    // make sure this goes out to stdout by not using Logger
    console.info(help);
    process.exit(2);
}

const portIdx = process.argv.indexOf('--port');
if (portIdx >= 0) {
    if (process.argv.length > portIdx + 1) {
        appConfig.port = process.argv[portIdx + 1];
    } else {
        console.error('missing port number or named pipe for --port option');
        process.exit(1);
    }
}

appConfig.mode = 'development';
appConfig.logLevel = 'debug';

////////////////////////////////////////////////////////////////////////////////
// Initialization

// make sure it exists
mkdirp.sync(appConfig.userDataPath);

// start the web server and web app
require(paths.src.www);

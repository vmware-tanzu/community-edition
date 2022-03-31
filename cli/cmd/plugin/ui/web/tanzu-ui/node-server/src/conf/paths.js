'use strict';

/* globals __dirname */

const path = require('path');

let serverDir = path.normalize(`${__dirname}/../..`);
let clientDistDir = path.normalize(`${serverDir}/../dist`);
let srcDir = `${serverDir}/src`;
let jsonDir = `${serverDir}/json`;
let mockLogsDir = `${serverDir}/mockLogs`;
let resourcesDir = `${serverDir}/resources`;
let supportDir = `${serverDir}/support`;

let paths = {
    // server/** directory paths
    directories: {
        clientDistDir,
        jsonDir,
        serverDir,
        srcDir,
        resourcesDir,
        supportDir
    },

    json: {
        mockJsonResults: `${jsonDir}/json-mock-response.json`
    },

    mockLogs: {
        deployment: `${mockLogsDir}/deployment.txt`
    },

    resources: {
        initErrorHtml: `${resourcesDir}/serverInitError.html`,
        mockTemplateFolder: `${resourcesDir}`
    },

    // server/src/** files
    src: {
        app: `${srcDir}/app`,
        appConfig: `${srcDir}/conf/appConfig`,
        bodyParser: `${srcDir}/services/bodyParser`,
        restApiRoutes: `${srcDir}/routes/api/restApiEndpoints`,
        util: `${srcDir}/services/util`,
        www: `${srcDir}/www`
    }
};

module.exports = paths;

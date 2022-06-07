// Library imports
const path = require('path');

const serverDir = path.normalize(`${__dirname}/../..`);
const clientDistDir = path.normalize(`${serverDir}/../dist`);
const srcDir = `${serverDir}/src`;
const jsonDir = `${serverDir}/json`;
const mockLogsDir = `${serverDir}/mockLogs`;
const resourcesDir = `${serverDir}/resources`;
const supportDir = `${serverDir}/support`;

const paths = {
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
        routes: {
            common: `${srcDir}/routes/api/endpoints/common`,
            docker: `${srcDir}/routes/api/endpoints/docker`,
            vsphere: `${srcDir}/routes/api/endpoints/vsphere`,
            aws: `${srcDir}/routes/api/endpoints/aws`,
            azure: `${srcDir}/routes/api/endpoints/azure`,
            management: `${srcDir}/routes/api/endpoints/management`,
            unmanaged: `${srcDir}/routes/api/endpoints/unmanaged`,
        },
        util: `${srcDir}/services/util`,
        www: `${srcDir}/www`
    }
};

module.exports = paths;

// Library imports
const express = require('express');
const winston = require('winston');

// App imports
const paths = require('../../../conf/paths');
const util = require(paths.src.util);

let mockDockerDaemonCounter = 0;
// eslint-disable-next-line new-cap
const router = express.Router({
    // '/Foo' different from '/foo'
    caseSensitive: true,
    // '/foo' and '/foo/' treated the same
    strict: false,
});

const readFile = util.readJsonFileWrapper(`${__dirname}/../responses/`);

router.get('/daemon', (req, res) => {
    winston.info('Mock UI VALIDATE DOCKER DAEMON');
    mockDockerDaemonCounter++;
    res.status(200);
    res.json({
        status: mockDockerDaemonCounter > 1 ? true : false,
    });
});

/**
 * Mock route for create docker cluster
 */
router.post('/create', (req, res) => {
    winston.info('Mock UI CREATE docker CLUSTER');
    res.status(200);
    res.json({});
});

/**
 * Mock route for apply config docker cluster
 */
router.post('/tkgconfig', (req, res) => {
    winston.info('Mock UI APPLY CONFIG');
    res.status(200);
    res.json({
        path: '/path/to/config',
    });
});

/**
 * Mock route for getting config file
 */
router.post('/config/export', (req, res) => {
    winston.info('Mock UI export config');
    res.status(200);
    res.json(readFile('provider-docker-export.json'));
});
/**
 * Import config
 */
router.post('/config/import', (req, res) => {
    winston.info('Mock UI IMPORT');
    res.status(200);
    res.json(readFile('provider-docker-import.json'));
});

module.exports = router;

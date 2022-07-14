// Library imports
const express = require('express');
const winston = require('winston');

// App imports
const paths = require('../../../conf/paths');
const util = require(paths.src.util);

// eslint-disable-next-line new-cap
const router = express.Router({
    // '/Foo' different from '/foo'
    caseSensitive: true,
    // '/foo' and '/foo/' treated the same
    strict: false
});

const readFile = util.readJsonFileWrapper(`${__dirname}/../responses/`);

/**
 * Mock route for getting unmanaged clusters
 */
router.get('/', (req, res) => {
    winston.info('Mock UI FETCH UNMANAGED CLUSTERS');
    res.status(200);
    res.json(readFile('unmanaged.json'));
});

/**
 * Mock route for create unmanaged cluster
 */
router.post('/', (req, res) => {
    winston.info('Mock UI CREATE UNMANAGED CLUSTER');
    res.status(200);
    res.json({});
});

module.exports = router;

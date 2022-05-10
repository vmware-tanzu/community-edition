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

router.get('/', (req, res) => {
    winston.info('Mock UI FETCH MANAGEMENT CLUSTERS');
    res.status(200);
    res.json(readFile('management-clusters.json'));
});

router.get('/:mcid/clusterclass/:ccid', (req, res) => {
    winston.info(`Mock UI FETCH CC: MCID=${req.params.mcid}, CCID=${req.params.ccid}`)
    res.status(200);
    res.json(readFile(req.params.ccid))
});

module.exports = router;

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
    strict: false,
});

const readFile = util.readJsonFileWrapper(`${__dirname}/../responses/`);

/**
 * Retrieve Azure account params from ENV variables
 */
router.get('/', (req, res) => {
    winston.info('Mock UI FETCH AZURE CREDENTIALS');
    res.status(200);
    res.json(readFile('provider-azure-credentials.json'));
});

/**
 * Verify Azure account credentials
 */
router.post('/', (req, res) => {
    winston.info('Mock UI VERIFY AZURE CREDENTIALS');
    if (req.body.tenantId && req.body.clientId && req.body.clientSecret && req.body.subscriptionId) {
        res.status(201);
        res.json({});
    } else {
        res.status(400);
        res.json({ message: 'Incorrect credentials' });
    }
});

/**
 * Retrieve Azure resource groups
 */
router.get('/resourcegroups', (req, res) => {
    winston.info('Mock UI FETCH AZURE RESOURCE GROUPS');
    if (req.query.location) {
        res.status(200);
        res.json(readFile('provider-azure-resourcegroups.json'));
    } else {
        res.status(400);
        res.json({ message: 'Missing resource group "region"' });
    }
});

/**
 * Mock route for importing config file
 */
router.post('/config/import', (req, res) => {
    winston.info('Mock UI import config');
    res.status(200);
    res.status(200);
    res.json(readFile('provider-azure-import.json'));
});

/**
 * Create an Azure resource group
 */
router.post('/resourcegroups', (req, res) => {
    winston.info('Mock UI CREATE AZURE RESOURCE GROUP');
    if (req.body.location && req.body.name) {
        res.status(201);
        res.json({});
    } else {
        res.status(400);
        res.json({ message: 'Missing either resource group "region" or "name"' });
    }
});

/**
 * Retrieve Azure VNets for a particular resource group
 */
router.get('/resourcegroups/:rgn/vnets', (req, res) => {
    winston.info('Mock UI RETRIEVE AZURE VNETS');
    if (req.params.rgn) {
        res.status(200);
        res.json(readFile('provider-azure-vnets.json'));
    } else {
        res.status(400);
        res.json({ message: 'Missing either resource group name' });
    }
});

/**
 * Retrieve Azure regions
 */
router.get('/regions', (req, res) => {
    winston.info('Mock UI RETRIEVE AZURE REGIONS');
    res.status(200);
    res.json(readFile('provider-azure-regions.json'));
});

router.get('/regions/:location/instanceTypes', (req, res) => {
    winston.info('Mock UI RETRIEVE AZURE REGIONS');
    res.status(200);
    res.json(readFile('provider-azure-instancetypes.json'));
});

/**
 * Retrieve os image
 */
router.get('/osimages', (req, res) => {
    winston.info('Mock UI RETRIEVE AZURE OS IMAGES');
    res.status(200);
    res.json(readFile('provider-azure-osimages.json'));
});

/**
 * Mock route for create Azure cluster
 */
router.post('/create', (req, res) => {
    winston.info('Mock UI CREATE AZURE CLUSTER');
    res.status(200);
    res.json({});
});

/**
 * Mock route for apply tkg config azure cluster
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
    res.status(200);
    res.json('Pretend this is a beautiful config file');
});
module.exports = router;

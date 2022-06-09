// Library imports
const express = require('express');
const winston = require('winston');

// App imports
const paths = require('../../../conf/paths');
const util = require(paths.src.util);

let mockvcNetworkRequestCounter = 0;
let mockResourcePoolRequestCounter = 0;
let mockOsImageRequestCounter = 0;

// eslint-disable-next-line new-cap
const router = express.Router({
    // '/Foo' different from '/foo'
    caseSensitive: true,
    // '/foo' and '/foo/' treated the same
    strict: false,
});

const readFile = util.readJsonFileWrapper(`${__dirname}/../responses/`);

/**
 * Mock route for apply tkg config vsphere cluster
 */
router.post('/tkgconfig', (req, res) => {
    winston.info('Mock UI APPLY CONFIG');
    res.status(200);
    res.json({
        path: '/path/to/config',
    });
});

/**
 * Mock route for create vsphere cluster
 */
router.post('/create', (req, res) => {
    winston.info('Mock UI CREATE VSPHERE CLUSTER');
    res.status(200);
    res.json({});
});

/**
 * Retrieve compute resources
 */
router.get('/compute', (req, res) => {
    winston.info('Mock UI RETRIEVE COMPUTE RESOURCES');
    res.status(200);
    res.json(readFile('provider-vsphere-compute.json'));
});

/**
 * Mock route for connect VC server
 */
router.post('/', (req, res) => {
    winston.info('Mock UI CONNECT VC API');
    if (
        (req.body.host === 'vsphere.local' || req.body.host === '1.1.1.1' || req.body.host === '2001:0db8:85a3:0000:0000:8a2e:0370:7334') &&
        (req.body.username === 'admin' || req.body.username === 'administrator') &&
        req.body.password === 'password'
    ) {
        res.status(200);
        res.json(readFile('provider-vsphere-connect.json'));
    } else if (req.body.host === '1.2.3.45') {
        res.status(500);
        res.json({ message: 'server error' });
    } else {
        res.status(403);
        res.json({ message: 'incorrect username or password' });
    }
});

/**
 * Mock route for getting VC datacenter network
 */
router.get('/networks', (req, res) => {
    winston.info('Mock UI VC NETWORK API');

    let vcNetworksResponse = readFile('provider-vsphere-networks.json');

    if (mockvcNetworkRequestCounter > 0) {
        vcNetworksResponse = readFile('provider-vsphere-networks-1.json');
    }

    mockvcNetworkRequestCounter++;

    res.status(200);
    res.json(vcNetworksResponse);
});

/**
 * Mock route for thumbsprint
 */

router.get('/thumbprint', (req, res) => {
    winston.info('Mock UI FETCH THUMBPRINT');
    res.status(200);
    res.json(readFile('provider-vsphere-thumbprint.json'));
});

/**
 * Mock route for getting VC datacenters
 */
router.get('/datacenters', (req, res) => {
    winston.info('Mock UI FETCH DATACENTERS');
    res.status(200);
    res.json(readFile('provider-vsphere-datacenters.json'));
});

/**
 * Mock route for getting datastores
 */
router.get('/datastores', (req, res) => {
    winston.info('Mock UI FETCH DATASTORES');
    res.status(200);
    res.json(readFile('provider-vsphere-datastores.json'));
});

/**
 * Mock route for getting vm folders
 */
router.get('/folders', (req, res) => {
    winston.info('Mock UI FETCH VM FOLDERS');
    res.status(200);
    if (req.body.dc === 'dc-1') {
        res.json(readFile('provider-vsphere-folders-1.json'));
    } else {
        res.json(readFile('provider-vsphere-folders-2.json'));
    }
});

/**
 * Mock route for getting VC datacenters
 */
router.get('/resourcepools', (req, res) => {
    winston.info('Mock UI FETCH RESOURCE POOLS');

    let resourcePoolsResponse = readFile('provider-vsphere-resourcepools.json');

    if (mockResourcePoolRequestCounter > 0) {
        resourcePoolsResponse = readFile('provider-vsphere-resourcepools-1.json');
    }

    mockResourcePoolRequestCounter++;

    setTimeout((_) => {
        res.status(200);
        res.json(resourcePoolsResponse);
    }, 5000);
});

/**
 * Mock route for getting VC os images
 */
router.get('/osimages', (req, res) => {
    winston.info('Mock UI FETCH DATACENTERS (#' + mockOsImageRequestCounter + ')');
    let osImageResponse = [];
    if (mockOsImageRequestCounter > 0) {
        osImageResponse = readFile('provider-vsphere-osimages.json');
    }

    mockOsImageRequestCounter++;

    res.status(200);
    res.json(osImageResponse);
});

router.post('/config/import', (req, res) => {
    winston.info('Mock UI IMPORT VC API');
    res.status(200);
    res.json('provider-vsphere-import.json');
});

/**
 * Mock route for getting config file
 */
router.post('/config/export', (req, res) => {
    winston.info('Mock UI export config');
    res.status(200);
    res.json('Pretend this is a beautiful config file');
});

module.exports = router;

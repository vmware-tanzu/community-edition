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
 * Retrieve AWS account params from ENV variables
 */
router.get('/', (req, res) => {
    winston.info('Mock UI FETCH AWS CREDENTIALS');
    const credentials = Math.random() > 0.2 ? readFile('provider-aws-credentials.json') : {};
    res.status(200);
    res.json(credentials);
});

/**
 * Verify AWS account credentials
 */
router.post('/', (req, res) => {
    winston.info('Mock UI VERIFY AWS CREDENTIALS');
    res.status(201);
    res.json({});
});

/**
 * Retrieve AWS vpc's
 */
router.get('/vpc', (req, res) => {
    winston.info('Mock UI RETRIEVE AWS VPCS');
    res.status(200);
    res.json(readFile('provider-aws-vpc.json'));
});

/**
 * Retrieve AWS availability zones
 */
router.get('/AvailabilityZones', (req, res) => {
    winston.info('Mock UI RETRIEVE AWS AVAILABILITY ZONES');
    res.status(200);
    res.json(readFile('provider-aws-az.json'));
});

/**
 * Retrieve AWS VPC CIDRS
 */
router.get('/subnets', (req, res) => {
    winston.info('Mock UI RETRIEVE AWS VPC SUBNET INFO');
    res.status(200);
    res.json(readFile('provider-aws-subnets.json'));
});

/**
 * Retrieve AWS node types
 */
router.get('/nodetypes', (req, res) => {
    winston.info('Mock UI RETRIEVE AWS NODE TYPES');
    res.status(200);
    res.json(readFile('provider-aws-nodetypes.json'));
});

/**
 * Retrieve AWS regions
 */
router.get('/regions', (req, res) => {
    winston.info('Mock UI RETRIEVE AWS REGIONS');
    res.status(200);
    res.json(readFile('provider-aws-regions.json'));
});

/**
 * Retrieve AWS profiles
 */
router.get('/profiles', (req, res) => {
    winston.info('Mock UI RETRIEVE AWS PROFILES');
    res.status(200);
    res.json(readFile('provider-aws-profiles.json'));
});

/**
 * Retrieve os image
 */
router.get('/osimages', (req, res) => {
    winston.info('Mock UI RETRIEVE AWS OS IMAGES');
    res.status(200);
    res.json(readFile('provider-aws-osimages.json'));
});

/**
 * Mock route for create aws cluster
 */
router.post('/create', (req, res) => {
    winston.info('Mock UI CREATE AWS CLUSTER');
    res.status(200);
    res.json({});
});

/**
 * Mock route for apply the config aws cluster
 */
router.post('/tkgconfig', (req, res) => {
    winston.info('Mock UI APPLY CONFIG');
    res.status(200);
    res.json({
        path: '/path/to/config'
    });
});

/**
 * Mock route for getting config file
 */
router.post('/config/export', (req, res) => {
    winston.info('Mock UI export config');
    res.status(200);
    res.json('Pretend this is a beautiful config file');
});

router.post('/config/import', (req, res) => {
    winston.info('Mock UI import config');
    res.status(200);
    res.json(readFile('provider-aws-import.json'));
});

module.exports = router;

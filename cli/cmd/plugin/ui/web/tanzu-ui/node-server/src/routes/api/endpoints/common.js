// Library imports
const express = require('express');
const winston = require('winston');

// App imports
const paths = require('../../../conf/paths');
const util = require(paths.src.util);
const appConfig = require(paths.src.appConfig);

// path to json file based responses
const mockJsonResultsRes = util.readJsonFileSync(paths.json.mockJsonResults, 'utf8');

// eslint-disable-next-line new-cap
const router = express.Router({
    // '/Foo' different from '/foo'
    caseSensitive: true,
    // '/foo' and '/foo/' treated the same
    strict: false
});

const readFile = util.readJsonFileWrapper(`${__dirname}/../responses/`);

/*
 * Mock route, for GET JSON file specific to design type
 */
router.get('/test/:urlparam', (req, res) => {
    winston.info('Mock UI GET API with URL param');
    res.status(200);
    res.json(mockJsonResultsRes);
});

/*
 * Mock route, for GET all JSON file
 */
router.get('/test', (req, res) => {
    winston.info('Mock UI GET API no URL param');
    res.status(200);
    res.json(mockJsonResultsRes);
});

/*
 * Mock route, for GET specified provider
 */
router.get('/providers', (req, res) => {
    winston.info('Mock UI GET PROVIDERS API');
    res.status(200);
    res.json(readFile('common-providers.json'));
});

/*
 * Mock route, for GET feature flags
 */
router.get('/features', (req, res) => {
    winston.info('Mock UI GET FEATURES API');
    res.status(200);
    res.json(readFile('common-features.json'));
});

/*
 * Mock route, for GET edition
 */
router.get('/edition', (req, res) => {
    winston.info('Mock UI GET EDITION API: ' + appConfig.edition);
    res.status(200);
    res.json(appConfig.edition);
});

/**
 * Mock route for connect Avi controller
 */
router.post('/avi', (req, res) => {
    winston.info('Mock UI CONNECT AVI CONTROLLER API');
    if ((req.body.host === 'avi.local') &&
        (req.body.username === 'administrator' || req.body.username === 'admin') &&
        (req.body.password === 'password')) {
        res.status(200);
        res.json({});
    } else {
        res.status(403);
        res.json({ message: 'incorrect username or password' });
    }
});

/**
 * Mock route for getting Avi clouds
 */
router.get('/avi/clouds', (req, res) => {
    winston.info('Mock UI FETCH AVI CLOUDS');
    res.status(200);
    res.json(readFile('common-avi-clouds.json'));
});

/**
 * Mock route for getting Avi service engine groups
 */
router.get('/avi/serviceenginegroups', (req, res) => {
    winston.info('Mock UI FETCH AVI SERVICE ENGINE GROUPS');
    res.status(200);
    res.json(readFile('common-avi-serviceenginegroups.json'));
});

/**
 * Mock route for getting VIP networks
 */
router.get('/avi/vipnetworks', (req, res) => {
    winston.info('Mock UI FETCH AVI VIP NETWORKS');
    res.status(200);
    res.json(readFile('common-avi-vipnetworks.json'));
});

/**
 * LDAP verification mock services
 */
router.post('/ldap/connect', (req, res) => {
    winston.info('Mock UI VERIFY LDAP CONNECTION');
    res.status(200);
});

router.post('/ldap/bind', (req, res) => {
    winston.info('Mock UI VERIFY LDAP BIND');
    res.status(200);
});

router.post('/ldap/users/search', (req, res) => {
    winston.info('Mock UI VERIFY LDAP USER SEARCH');
    res.status(200);
});

router.post('/ldap/groups/search', (req, res) => {
    winston.info('Mock UI VERIFY LDAP GROUP SEARCH');
    res.status(200);
});

router.post('/ldap/disconnect', (req, res) => {
    winston.info('Mock UI VERIFY LDAP DISCONNECT');
    res.status(200);
});

module.exports = router;

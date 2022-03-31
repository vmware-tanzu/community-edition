'use strict';

const paths = require('../../conf/paths');
const util = require(paths.src.util);
const express = require('express');
const appConfig = require(paths.src.appConfig);
const winston = require('winston');
const ENDPOINT = appConfig.apiEndpoint;

// path to json file based responses
let mockJsonResultsRes = util.readJsonFileSync(paths.json.mockJsonResults, 'utf8');
let mockDockerDaemonCounter = 0;

// eslint-disable-next-line new-cap
let router = express.Router({
    // '/Foo' different from '/foo'
    caseSensitive: true,
    // '/foo' and '/foo/' treated the same
    strict: false
});

/*
 * Mock route, for GET JSON file specific to design type
 */
router.get(`${ENDPOINT}/test/:urlparam`, (req, res) => {
    winston.info('Mock UI GET API with URL param');
    res.status(200);
    res.json(mockJsonResultsRes);
});

/*
 * Mock route, for GET all JSON file
 */
router.get(`${ENDPOINT}/test`, (req, res) => {
    winston.info('Mock UI GET API no URL param');
    res.status(200);
    res.json(mockJsonResultsRes);
});

/**
 * Create an Azure resource group
 */
router.post(`${ENDPOINT}/test`, (req, res) => {
    winston.info('Mock UI POST API body param region');
    if (req.body.region) {
        res.status(201);
        res.json({});
    } else {
        res.status(400);
        res.json({ message: "Missing mock 'region'" });
    }
});

/*********************************   DOCKER   **********************************/

router.get(`${ENDPOINT}/providers/docker/daemon`, (req, res) => {
    winston.info('Mock TCE UI VALIDATE DOCKER DAEMON');
    mockDockerDaemonCounter++;
    res.status(200);
    res.json(
        {
            status: mockDockerDaemonCounter > 1 ? true : false
        }
    );
});

/**
 * Mock route for create docker cluster
 */
 router.post(`${ENDPOINT}/providers/docker/create`, (req, res) => {
    winston.info('Mock TCE UI CREATE docker CLUSTER');
    res.status(200);
    res.json({});
});

/**
 * Mock route for apply tce config docker cluster
 */
router.post(`${ENDPOINT}/providers/docker/tceconfig`, (req, res) => {
    winston.info('Mock TCE UI APPLY TCE CONFIG');
    res.status(200);
    res.json({
        path: "/path/to/config"
    });
});


module.exports = router;

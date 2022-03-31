'use strict';

//
// Web Application Configuration
//
const _ = require('lodash');
const assert = require('assert');
const path = require('path');
const paths = require('./paths');
const libUtil = require(paths.src.util);

//
// Module
//

// @see #def.port
let __port;
// @see #def.userDataPath
let __userDataPath;
// @see #def.logPath
let __logPath;
// @see #def.clientPath
let __clientPath;
// @see #def.logLevel
let __logLevel;
// @see #def.apiEndpoint
let __apiEndpoint;
// @see #def.sampleJsonDirectory
let __sampleJsonDirectory;
// @see #def.__serverErrorState
let __serverErrorState;

// [internal]
// @const {Number} Default app port.
const DEFAULT_PORT = 8008;
// @const {String} Default log level
const DEFAULT_LOG_LEVEL = 'info';

// {Object.<namespace:String,Object>} module definition
let def = {

    /**
     * Full path to the user data directory where user-specific files can be stored.
     * @readonly
     * @type {String}
     */
    get userDataPath() {
        return __userDataPath;
    },

    /**
     * Full path to the directory where log files should be located.
     * @readonly
     * @type {String}
     */
    get logPath() {
        return __logPath;
    },

    /**
     * Root API path. Always starts with a slash, never ends with one.
     * @readonly
     * @type {String}
     */
    get apiEndpoint() {
        return __apiEndpoint;
    },

    /**
     * Minimum web server port.
     * @readonly
     * @type {Number}
     */
    get minPort() {
        return 8000;
    },

    /**
     * Web server port or named pipe. Defaults to port 8008.
     *
     * If a number is specified, it must be at least `minPort` (and will be adjusted as
     *  such if it isn't). If a string is specified, it's expected to be a named pipe.
     *
     * @type {Number|String}
     */
    get port() {
        return __port;
    },
    set port(newValue) {
        if (_.isString(newValue) && newValue) {
            // named pipe
            __port = newValue;
        } else {
            // assume a port
            // use default if NaN or 0
            __port = parseInt(newValue, 10) || DEFAULT_PORT;
            // ensure at least minimum
            __port = Math.max(__port, this.minPort);
        }
    },

    /**
     * Web application log level.
     * @type {Number.<vmw.iup.log.Logger.level>}
     */
    get logLevel() {
        return __logLevel;
    },
    set logLevel(newValue) {
        __logLevel = newValue;
    },

    /**
     * Absolute path to location of client web app to serve statically.
     * @type {String|undefined}
     */
    get clientPath() {
        return __clientPath;
    },
    set clientPath(newValue) {
        assert(newValue && _.isString(newValue), 'Client Path must be a string');

        __clientPath = newValue;
    },

    /**
     * Absolute path to location of folder containing sampleJson
     * @type {String|undefined}
     */
    get sampleJsonDirectory() {
        return __sampleJsonDirectory;
    },
    set sampleJsonDirectory(newValue) {
        assert(newValue && _.isString(newValue), 'Json path must be a string');
        __sampleJsonDirectory = newValue;
    },

    /**
     * Address of mock server (for development mode)
     * @Type {Boolean}
     */
    get serverErrorState() {
        return __serverErrorState;
    },
    set serverErrorState(newValue) {
        assert(!_.isUndefined(newValue) &&
            _.isBoolean(newValue), 'Server error state must be a boolean');
        __serverErrorState = newValue;
    },

    // [override]
    toString() {
        return libUtil.stringify(this, {
            typeName: 'server.appConfig',
            excludeFunctions: true
        });
    }
};

//// initialize
__userDataPath = `${paths.directories.serverDir}/user_data`;
__logPath = path.join(__userDataPath, 'logs');
__port = DEFAULT_PORT;
__logLevel = DEFAULT_LOG_LEVEL;
__apiEndpoint = '/api';
__sampleJsonDirectory = `${__dirname}/../../json`;
__serverErrorState = false;

// set default mode
def.mode = 'development';

module.exports = def;

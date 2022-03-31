'use strict';

//
// General utility library.
// NOTE: This module must NOT depend on 'yllr'; exceptionally, it should use Node.assert
//

const _ = require('lodash');
const assert = require('assert');
const path = require('path');
const fs = require('fs');
const PropertiesReader = require('properties-reader');

//
// Functions
//
/**
 * Determines if a path points to a directory that exists.
 * @param {String} dir Path to check.
 * @returns {Boolean} `true` if `dir` is an existing directory; `false` otherwise.
 */
let isDirSync = function(dir) {
    return (fs.existsSync(dir) && fs.statSync(dir).isDirectory());
};

/**
 * Generates a string representation of a given object (typically a type or shape).
 * @param {Object} obj Object to stringify
 *  (can be any type of JSON-serializable value).
 * @param {Object} [options] Options for stringification:
 *  - excludeFunctions:Boolean [false] If truthy, functions are excluded from output.
 *  - [typeName:String] Optional type name to add to the representation. Used
 *      only if `obj` is an object (not an array, not a primitive).
 * @returns {String|undefined} String representation of the object; `undefined`
 *  if `obj === undefined`.
 */
let stringify = function(obj, options) {
    options = _.defaults(options, {
        excludeFunctions: false,
        typeName: undefined
    });

    // make sure arrays and objects come out fully-speced
    let rep = obj === undefined ? undefined : JSON.stringify(obj, (key, value) => {
        if (_.isFunction(value)) {
            // likely an iface custom validator
            return options.excludeFunctions ? undefined : '<function>';
        } else if (value === undefined) {
            return '<undefined>';
        }

        return value;
    });

    if (options.typeName && _.isObject(obj)) {
        // result of `JSON.stringify()` will be an object, as '{...}', so insert
        //  the type name after '{', resulting in '{<typeName> ...}'
        rep = rep[0] + options.typeName + ' ' + rep.substr(1);
    }

    return rep;
};

/**
 * Makes a filename generator function for use with the 'rotating-file-stream' module.
 * @param {String} filepath Path to the primary log file. What follows the last
 *  period will be used as the extension. Archived files will be placed in
 *  an 'older' subdir.
 * @returns {Function} A rotating-file-stream log filename generator.
 */
let makeRotatingFilenameGenerator = function(filepath) {
    assert(filepath && _.isString(filepath),
        'makeRotatingFilenameGenerator(): filepath must be a non-empty string: ' +
        filepath);

    const ARCHIVE_SUBPATH = 'older';
    let fileext = path.extname(filepath);
    let filename = path.basename(filepath, fileext);
    let dirpath = path.dirname(filepath);

    // remove leading '.'
    fileext = fileext.substr(1);

    // rotating file stream filename generator
    return function(time, idx) {
        if (!time) {
            // non-rotated filename; also the name of the current log file (i.e. latest
            //  log entries will always be written to this file, and older ones rotated
            //  out to files with timestamps below)
            return path.join(dirpath, filename + '.' + fileext);
        }

        // add timestamp and index to name of rotated-out file
        return path.join(dirpath, ARCHIVE_SUBPATH,
            filename + '-' +
            time.getFullYear() +
            _.padStart((time.getMonth() + 1).toString(10), 2, '0') +
            _.padStart(time.getDate().toString(10), 2, '0') + '+' +
            _.padStart(time.getHours().toString(10), 2, '0') +
            _.padStart(time.getMinutes().toString(10), 2, '0') + '+' +
            _.padStart(idx.toString(10), 2, '0') + '.' +
            fileext
        );
    };
};

/**
 * Read in a .properties file
 * @param {String} filepath Absolute or relative path to the properties file to be read
 * @param {Array<String>} properties properties to read
 * @returns {Object} An object corresponding to the parsed properties file,
 *  empty object if file does not exist
 */
let readPropertiesFileSync = function(filepath, properties) {
    let propertiesRes = PropertiesReader(filepath);
    return properties.map(property => propertiesRes.get(property));
};

/**
 * Read in a json file with specified encoding
 * @param {String} filepath Absolute or relative path to the json file to be read
 * @param {String} [encoding] Encoding to use when reading the file, utf8 by default
 * @returns {Object} An object corresponding to the parsed json file,
 *  empty object if file does not exist
 */
let readJsonFileSync = function(filepath, encoding) {
    assert(_.isString(filepath), 'readJsonFileSync(): filepath: expecting ' +
        'String, got ' + filepath);
    if (!fs.existsSync(filepath)) {
        return {};
    }

    if (typeof (encoding) === 'undefined') {
        encoding = 'utf8';
    }
    let file = fs.readFileSync(filepath, encoding);
    try {
        JSON.parse(file);
    } catch (e) {
        return file;
    }
    return JSON.parse(file);
};

// module definition
module.exports = {
    // methods
    isDirSync,
    stringify,
    makeRotatingFilenameGenerator,
    readJsonFileSync,
    readPropertiesFileSync
};

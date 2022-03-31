const bodyParser = require('body-parser');
const clonedBodyParser = Object.assign({}, bodyParser);
const _json = bodyParser.json;
const _urlEncoded = bodyParser.urlencoded;

const isMultipartRequest = function(req) {
    let contentTypeHeader = req.headers['content-type'];
    return contentTypeHeader && contentTypeHeader.indexOf('multipart') > -1;
};

/**
 * @method Override of bodyParser.json function
 * Returns a modified json middleware
 * @returns {Function}  Modified instance of json() middleware
 */
clonedBodyParser.json = function() {
    let _jsonInstance = _json(...arguments);

    /**
     * @inner
     * body parser does not handle multipart bodies, due to their complex
     * and typically large nature. Bypass body parsing middleware in case
     * of multipart requests
     *  @returns {Function}  Modified instance of json() middleware
     *  @param req
     *  @param res
     *  @param next
     */
    return (req, res, next) => {
        if (isMultipartRequest(req)) {
            return next();
        }
        return _jsonInstance(req, res, next);
    };
};

/**
 * @method Override of bodyParser.urlencoded function
 * Returns a modified json middleware
 * @returns {Function}  Modified instance of urlencoded() middleware
 */
clonedBodyParser.urlencoded = function() {
    let _urlEncodedInstance = _urlEncoded(...arguments);

    /**
     * @inner
     * body parser does not handle multipart bodies, due to their complex
     * and typically large nature. Bypass body parsing middleware in case
     * of multipart requests
     *  @returns {Function}  Modified instance of json() middleware
     *  @param req
     *  @param res
     *  @param next
     */
    return (req, res, next) => {
        if (isMultipartRequest(req)) {
            return next();
        }
        return _urlEncodedInstance(req, res, next);
    };
};

module.exports = clonedBodyParser;

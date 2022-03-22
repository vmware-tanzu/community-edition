'use strict'
const { execSync } = require("child_process")
const fs = require('fs')
const yaml = require('js-yaml')

// configPath should be:
// Darwin: os.homedir() + '/.config/tanzu/config.yaml'
function tanzuEdition(configPath) {
    try {
        let fileContents = fs.readFileSync(configPath, 'utf8')
        const data = yaml.load(fileContents)
        if (!data) {
            console.log('opened config, but unable to parse data')
        } else {
            const edition = data?.clientOptions?.cli?.edition
            if (!edition) {
                console.log('Config file at ' + configPath + ' does not seem to include clientOptions.cli.edition: ' + JSON.stringify(data))
            } else {
                return edition.toUpperCase()
            }
        }
    } catch (e) {
        console.log(e)
    }
    return ''
}

// "command" should be "which" on Darwin/Linux and "where" for Windows
function tanzuPath(command) {
    const tanzuCommand = command + ' tanzu'
    const stdio = '' + execSync(tanzuCommand)
    const parts = stdio.match(/(.*)\n/)
    let path = ''
    if (parts && parts.length > 1) {
        path = parts[1]
    } else {
        console.log('Unable to parse Tanzu path from output of "' + tanzuCommand + '" command: ' + stdio)
    }
    return path
}

function tanzuVersion() {
    const stdio = '' + execSync('tanzu version')
    const parts = stdio.match(/version: (.*)\n/)
    let version = ''
    if (parts && parts.length > 1) {
        version = parts[1]
    } else {
        console.log('Unable to parse Tanzu version from output of "tanzu version" command: ' + stdio)
    }
    return version
}

module.exports.tanzuEdition = tanzuEdition
module.exports.tanzuPath = tanzuPath
module.exports.tanzuVersion = tanzuVersion

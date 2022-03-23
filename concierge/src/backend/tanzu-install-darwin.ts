const { spawnSync } = require("child_process")
const os = require( 'os' )
const fs = require('fs')

const tanzuUtil = require('./tanzu-install-util')
import { ProgressMessage } from '../models/progressMessage'
import { ExistingInstall } from '../models/existingInstall'


function checkExistingInstallationDarwin(): ExistingInstall {
    const version = tanzuUtil.tanzuVersion()
    const path = tanzuUtil.tanzuPath('which')
    const configPath = os.homedir() + '/.config/tanzu/config.yaml'
    const edition = tanzuUtil.tanzuEdition(configPath)

    const result = {path, version, edition}
    console.log('Existing install: ' + JSON.stringify(result))
    return result
}

function installDarwin(progressMessenger) {
    console.log('Installing darwin...')
    progressMessenger({message: 'Starting installation on darwin'})
    const tarballDir = __dirname + '/../assets/tanzu-release'
    if (!darwinTestPath(tarballDir)) {
        progressMessenger({error: true, finished: true, message: 'ERROR: unfortunately, we are not able to find the directory where we expect an installation tarball: ' + tarballDir + '' +
            ', so we\'ll have to abandon the installation effort. So sorry.'} )
        return false
    }
    progressMessenger({step: 'Check prerequisites', message: 'Directory exists: ' + tarballDir})

    const tarballFile = 'tce-darwin-amd64-v0.10.0.tar.gz'
    const tarballFullPath = tarballDir + '/' + tarballFile
    if (!darwinTestPath(tarballFullPath)) {
        progressMessenger({step: 'Check prerequisites', error: true, finished: true, message: 'ERROR: unfortunately, we are not able to find the installation tarball: ' + tarballFullPath + '' +
            ', so we\'ll have to abandon the installation effort. So sorry.'} )
        return false
    }

    const untarResult = darwinUntar(tarballDir, tarballFile)
    const step = 'Unpack tanzu'
    if (untarResult.error) {
        progressMessenger({step, error: true, finished: true, message: 'ERROR: unfortunately, we encountered an error trying to untar the installation tarball: ' + tarballFullPath + '' +
            ', so we\'ll have to abandon the installation effort. So sorry.', details: untarResult.message + '\n' + untarResult.details} )
        console.log('ERROR during untar: ' + JSON.stringify(untarResult))
        return false
    }
    console.log('SUCCESS untar: ' + JSON.stringify(untarResult))
    progressMessenger({...untarResult, step})

    return true
}

function darwinExec(command, args) : ProgressMessage {
    const result = {message: '', details: '', error: false }
    try {
        const syncResult = spawnSync(command, args, {stdio: 'pipe', encoding: 'utf8'})
        result.message = syncResult.stdout.toString()
        result.details = syncResult.stderr.toString()
    } catch (e) {
        console.log(e)
        result.message = 'ERROR: ' + e.toString()
        result.error = true
    }
    return result
}

function darwinTestPath(path) {
    try {
        fs.accessSync(path)
        return true
    } catch (e) {
        console.log(e)
        return false
    }
}

function darwinUntar(tarballDir, tarballFile) : ProgressMessage {
    const result = darwinExec('tar', ['xzvf', tarballDir + '/' + tarballFile, '-C', tarballDir])
    if (!result.error) {
        result.message =  'Successful untar of ' + tarballDir + '/' + tarballFile
    }
    result.stepComplete = true
    return result
}

module.exports.install = installDarwin
module.exports.checkExistingInstallation = checkExistingInstallationDarwin

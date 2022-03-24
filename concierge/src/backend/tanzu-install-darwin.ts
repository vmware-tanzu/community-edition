const { spawnSync } = require("child_process")
const os = require( 'os' )
const fs = require('fs')

const tanzuUtil = require('./tanzu-install-util.ts')
import { ExistingInstall } from '../models/existingInstall'
import { ProgressMessage, ProgressMessenger } from '../models/progressMessage'


function checkExistingInstallationDarwin(): ExistingInstall {
    const version = tanzuUtil.tanzuVersion()
    const path = tanzuUtil.tanzuPath('which')
    const configPath = os.homedir() + '/.config/tanzu/config.yaml'
    const edition = tanzuUtil.tanzuEdition(configPath)

    const result = {path, version, edition}
    console.log('Existing install: ' + JSON.stringify(result))
    return result
}

const darwinSteps = [
    { name: 'Check prerequisites', execute: darwinTarballCheck },
    { name: 'Unpack tanzu', execute: darwinTarballUnpack }
]

const tarballDir = __dirname
const tarballFile = 'tce-darwin-amd64-v0.10.0.tar.gz'   // TODO: dynamically find the available tarball(s)
const tarballFullPath = tarballDir + '/' + tarballFile

function darwinTarballCheck(progressMessenger: ProgressMessenger) : boolean {
    console.log('Installing darwin...')
    progressMessenger.report({message: 'Here we go... (starting installation on Mac OS)'})

    if (!darwinTestPath(tarballDir)) {
        progressMessenger.report({error: true, stepComplete: true, message: 'ERROR: unfortunately, we are not able to find the directory where we expect an installation tarball: ' + tarballDir + '' +
                ', so we\'ll have to abandon the installation effort. So sorry.'} )
        return false
    }
    progressMessenger.report({step: 'Check prerequisites', message: 'Directory exists: ' + tarballDir})

    const tarballFullPath = tarballDir + '/' + tarballFile
    if (!darwinTestPath(tarballFullPath)) {
        progressMessenger.report({step: 'Check prerequisites', error: true, stepComplete: true, message: 'ERROR: unfortunately, we are not able to find the installation tarball: ' + tarballFullPath + '' +
                ', so we\'ll have to abandon the installation effort. So sorry.'} )
        return false
    }

    return true
}

function darwinTarballUnpack(progressMessenger: ProgressMessenger) : boolean {
    const untarResult = darwinUntar(tarballDir, tarballFile)
    const step = 'Unpack tanzu'
    if (untarResult.error) {
        progressMessenger.report({step, error: true, stepComplete: true, message: 'ERROR: unfortunately, we encountered an error trying to untar' +
                ' the installation tarball: ' + tarballFullPath + '' +
                ', so we\'ll have to abandon the installation effort. So sorry.', details: untarResult.message + '\n' + untarResult.details} )
        console.log('ERROR during untar: ' + JSON.stringify(untarResult))
        return false
    }
    console.log('SUCCESS untar: ' + JSON.stringify(untarResult))
    progressMessenger.report({...untarResult, step})

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

module.exports.steps = darwinSteps
module.exports.checkExistingInstallation = checkExistingInstallationDarwin

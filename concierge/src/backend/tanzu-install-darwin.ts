const { spawnSync } = require("child_process")
const os = require( 'os' )
const fs = require('fs')

const tanzuUtil = require('./tanzu-install-util.ts')
import { ExistingInstallation, InstallationState, InstallationTarball, InstallStep } from '../models/installation'
import { ProgressMessage, ProgressMessenger } from '../models/progressMessage'


function checkExistingInstallationDarwin(): ExistingInstallation {
    const version = tanzuUtil.tanzuVersion()
    const path = tanzuUtil.tanzuPath('which')
    const configPath = os.homedir() + '/.config/tanzu/config.yaml'
    const edition = tanzuUtil.tanzuEdition(configPath)

    const result = {path, version, edition}
    console.log('Existing install: ' + JSON.stringify(result))
    return result
}

const darwinSteps = [
    { name: 'Check prerequisites', execute: darwinTarballCheck } as InstallStep,
    { name: 'Unpack tanzu', execute: darwinTarballUnpack } as InstallStep,
]


function darwinTarballCheck(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    console.log('Installing darwin...')
    progressMessenger.report({message: 'Here we go... (starting installation on Mac OS)'})

    const tarballDir = __dirname
    const tarballFile = 'tce-darwin-amd64-v0.10.0.tar.gz'   // TODO: dynamically find the available tarball(s)

    if (!darwinTestPath(tarballDir)) {
        progressMessenger.report({
            error: true,
            stepComplete: true,
            message: 'ERROR: unfortunately, we are not able to find the directory where we expect an installation tarball: ' + tarballDir + '' +
                ', so we\'ll have to abandon the installation effort. So sorry.'
        })
        return {...state, stop: true}
    }
    progressMessenger.report({step: state.currentStep, message: 'Directory exists: ' + tarballDir})

    const tarballFullPath = tarballDir + '/' + tarballFile
    if (!darwinTestPath(tarballFullPath)) {
        progressMessenger.report({
            step: state.currentStep,
            error: true,
            stepComplete: true,
            message: 'ERROR: unfortunately, we are not able to find the installation tarball: ' + tarballFullPath + '' +
                ', so we\'ll have to abandon the installation effort. So sorry.'
        })
        return {...state, stop: true}
    }

    const tarballInfo = { dir: tarballDir, file: tarballFile, fullPath: tarballFullPath }
    return {...state, tarballInfo }
}

function darwinTarballUnpack(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.tarballInfo) {
        progressMessenger.report({step: state.currentStep, error: true, stepComplete: true,
            message: 'ERROR: unfortunately, some programmer decided we should try to UNPACK the tarball BEFORE detecting the ' +
                ' required information ABOUT the tarball! We\'ll have to abandon the installation effort. So sorry.'} )
        return {...state, stop: true}
    }
    const untarResult = darwinUntar(state.tarballInfo)
    if (untarResult.error) {
        progressMessenger.report({step: state.currentStep, error: true, stepComplete: true,
            message: 'ERROR: unfortunately, we encountered an error trying to untar the installation tarball: ' +
                state.tarballInfo.fullPath + ', so we\'ll have to abandon the installation effort. So sorry.',
            details: untarResult.message + '\n' + untarResult.details} )
        console.log('ERROR during untar: ' + JSON.stringify(untarResult))
        return {...state, stop: true}
    }
    console.log('SUCCESS untar: ' + JSON.stringify(untarResult))
    progressMessenger.report({...untarResult, step: state.currentStep})

    return state
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

function darwinUntar(tarballInfo: InstallationTarball) : ProgressMessage {
    const result = darwinExec('tar', ['xzvf', tarballInfo.fullPath, '-C', tarballInfo.dir])
    if (!result.error) {
        result.message =  'Successful untar of ' + tarballInfo.fullPath
    }
    result.stepComplete = true
    return result
}

module.exports.steps = darwinSteps
module.exports.checkExistingInstallation = checkExistingInstallationDarwin

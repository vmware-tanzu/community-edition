const { spawnSync } = require("child_process")
const os = require( 'os' )
const fs = require('fs')

const util = require('./tanzu-install-util.ts');
import { ExistingInstallation, InstallationState, InstallationTarball } from '../models/installation'
import { ProgressMessage, ProgressMessenger } from '../models/progressMessage'

const darwinSteps = [
    { name: 'Check prerequisites', execute: darwinTarballCheck },
    { name: 'Unpack tanzu', execute: darwinTarballUnpack },
    { name: 'Set data paths', execute: darwinSetDataPaths },
    { name: 'Set binary path', execute: darwinSetBinPath },
    { name: 'Delete existing Tanzu binary', execute: darwinDeleteExistingTanzuIfNec },
    { name: 'Copy new Tanzu binary', execute: darwinCopyTanzuBinary },
    { name: 'Copy uninstall script', execute: darwinCopyUninstallScript },
    { name: 'Install plugins', execute: darwinInstallPlugins },
    { name: 'Add repos', execute: darwinAddTanzuReposIfNec },
]

function checkExistingInstallationDarwin(): ExistingInstallation {
    const version = util.tanzuVersion()
    const path = util.tanzuPath('which')
    const configPath = os.homedir() + '/.config/tanzu/config.yaml'
    const edition = util.tanzuEdition(configPath)

    const result = {path, version, edition}
    console.log('Existing install: ' + JSON.stringify(result))
    return result
}

function darwinTarballCheck(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    console.log('Installing darwin...')
    progressMessenger.report({message: 'Here we go... (starting installation on Mac OS)'})

    const tarballDir = __dirname
    const tarballFile = 'tce-darwin-amd64-v0.10.0.tar.gz'   // TODO: dynamically find the available tarball(s)

    if (!darwinTestPath(tarballDir)) {
        const message = 'ERROR: unfortunately, we are not able to find the directory where we expect an installation tarball: ' + tarballDir + '' +
        ', so we\'ll have to abandon the installation effort. So sorry.'
        return util.reportError(message, progressMessenger, state)
    }
    progressMessenger.report({step: state.currentStep, message: 'Directory exists: ' + tarballDir})

    const tarballFullPath = tarballDir + '/' + tarballFile
    if (!darwinTestPath(tarballFullPath)) {
        const message = 'ERROR: unfortunately, we are not able to find the installation tarball: ' + tarballFullPath + '' +
            ', so we\'ll have to abandon the installation effort. So sorry.'
        return util.reportError(message, progressMessenger, state)
    }

    const newState = {...state, tarballInfo: { dir: tarballDir, file: tarballFile, fullPath: tarballFullPath } }
    return util.reportMessage('Found tarball ' + tarballFullPath, true, progressMessenger, newState)
}

function darwinTarballUnpack(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.tarballInfo) {
        return util.reportMissingPrerequisite('detect tarball information', progressMessenger, state)
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

function darwinSetDataPaths(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    const pathTanzuDataHome = os.homedir() + '/.local/share'
    const pathTanzuConfigHome = os.homedir() + '/.config'

    const newState = {...state, pathTanzuConfigHome, pathTanzuDataHome}
    return util.reportMessage('Set data paths', true, progressMessenger, newState)
}

function darwinSetBinPath(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    const resultEnvPath = darwinExec('echo', ['$PATH'])
    if (resultEnvPath.error) {
        progressMessenger.report(resultEnvPath)
        return {...state, stop: true}
    }
    const envPath = resultEnvPath.message
    const preferredPath = os.homedir() + '/bin:'
    const defaultPath = '/usr/local/bin'
    const pathTanzuInstallation = darwinEnvPathContains(envPath, preferredPath) ? preferredPath : defaultPath
    const newState = {...state, pathTanzuInstallation}
    return util.reportDetails('Set binary path', 'Binary path=' + pathTanzuInstallation, true, progressMessenger, newState)
}

function darwinDeleteExistingTanzuIfNec(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.existingInstallation) {
        return util.reportMessage('No existing tanzu binary found (so skipping delete)', true, progressMessenger, state)
    }
    const result = darwinExec('rm', ['-f', state.existingInstallation.path])
    if (result.error) {
        return util.reportError(result.message, progressMessenger, state)
    }
    const message =  'Successful removal of ' + state.existingInstallation.path
    return util.reportMessage(message, true, progressMessenger, state)
}

function darwinCopyTanzuBinary(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.pathInstallFiles) {
        return util.reportMissingPrerequisite('set the source binary directory', progressMessenger, state)
    }
    if (!state.pathTanzuInstallation) {
        return util.reportMissingPrerequisite('set the target binary directory', progressMessenger, state)
    }

    const pathSourceBinary = state.pathInstallFiles + '/tanzu'
    const pathTargetBinary = state.pathTanzuInstallation + '/tanzu'
    const result = darwinExec('install', [pathSourceBinary, pathTargetBinary])
    if (result.error) {
        return util.reportError(result.message, progressMessenger, state)
    }
    const message =  'Successful copy of tanzu binary to ' + pathTargetBinary
    return util.reportMessage(message, true, progressMessenger, state)
}

function darwinCopyUninstallScript(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.pathInstallFiles) {
        return util.reportMissingPrerequisite('set the source binary directory', progressMessenger, state)
    }
    if (!state.pathTanzuDataHome) {
        return util.reportMissingPrerequisite('set the tanzu data directory', progressMessenger, state)
    }
    const tceDataDir = state.pathTanzuDataHome + '/tce'
    const resultMakeTceDir = darwinExec('mkdir', ['-p', tceDataDir])
    if (resultMakeTceDir.error) {
        return util.reportError(resultMakeTceDir.message, progressMessenger, state)
    }
    // TODO: log detail of creating/ensuring tce dir

    const sourceScript = state.pathInstallFiles + '/uninstall.sh'
    const resultInstallScript = darwinExec('install', [sourceScript, tceDataDir])
    if (resultInstallScript.error) {
        return util.reportError(resultInstallScript.message, progressMessenger, state)
    }
    const message = 'copied uninstall script to ' + tceDataDir
    return util.reportMessage(message, true, progressMessenger, state)
}

function darwinInstallPlugins(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    // TODO: implement plugin installation
    const message = 'PRETENDED TO INSTALL TANZU PLUGINS'
    return util.reportMessage(message, true, progressMessenger, state)
}

function darwinAddTanzuReposIfNec(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    // TODO: implement adding TCE repos
    const message = 'PRETENDED TO ADD TANZU REPO'
    return util.reportMessage(message, true, progressMessenger, state)
}

//============================================================
// helper methods
//============================================================
function darwinExec(command: string, args: string[]) : ProgressMessage {
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

function darwinEnvPathContains(envPath, target) {
    return util.stringContains(':' + envPath + ':', ':' + target + ':')
}

module.exports.steps = darwinSteps
module.exports.checkExistingInstallation = checkExistingInstallationDarwin

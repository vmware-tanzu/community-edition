'use strict'

import { ExistingInstallation, InstallationArchive, InstallationState, InstallData } from '../models/installation';
import { ProgressMessage, ProgressMessenger } from '_models/progressMessage';
import os from 'os';
const utilExec = require('./tanzu-exec-util.ts');
const utilInstall = require('./tanzu-install-util.ts');

function preinstallWin32(): ExistingInstallation {
    return undefined
}

const win32Steps = [
    { name: 'Check prerequisites', execute: utilInstall.checkInstallationArchive },
    { name: 'Set data paths', execute: winSetDataPaths },
    { name: 'Unpack tanzu', execute: winUnzipArchive },
]

function winUnzipArchive(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    return utilInstall.unpackArchive(winUnzip, state, progressMessenger)
}

function winUnzip(archive: InstallationArchive, dstDir: string) : ProgressMessage {
    const result = utilExec.exec({}, 'Expand-Archive', '-Path', archive.fullPath, '-DestinationPath', dstDir)
    if (!result.error) {
        result.message =  'Successful unzip of ' + archive.fullPath
    }
    result.stepComplete = true
    return result
}

function winSetDataPaths(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    // TODO: verify these paths are correct
    const userProfile = os.homedir()
    const dirTanzuData = os.homedir() + '/.local/share'
    const pathTanzuConfig = userProfile + '\\.config\\tanzu\\config.yaml'
    const dirTanzuConfig = userProfile + '\\.config'
    const dirTanzuTmp = userProfile + '\\.tanzu\\concierge\\tmp'

    const tmpResult = winEnsureDir(dirTanzuTmp)
    if (tmpResult.error) {
        const message = `Unable to ensure tmp dir (for unzipping) ${dirTanzuTmp}`
        console.log(`ERROR: ${message}; raw command result ${JSON.stringify(tmpResult)}`)
        return utilInstall.reportError(message, progressMessenger, state)
    }

    const newState = {...state, pathTanzuConfig, dirTanzuData, dirTanzuConfig, dirTanzuTmp}
    return utilInstall.reportStepComplete('Set data paths', progressMessenger, newState)
}

function winEnsureDir(dir: string): ProgressMessage {
    return utilExec.exec({},'mkdir', dir)
}


const installData = {
    steps: win32Steps,
    msgStart: 'Here we go... (starting installation on Windows OS)',
    msgFailed: 'So sorry the installation did not succeed. It\'s probably Bill Gates\' fault.' +
        ' Please try again after any known issues are addressed.',
    msgSucceeded: 'You\'re ready to start using Tanzu!',
} as InstallData
module.exports.installData = installData
module.exports.preinstall = preinstallWin32

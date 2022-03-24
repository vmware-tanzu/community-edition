'use strict'
import { InstallStep } from '../models/installStep';
import { ProgressMessenger } from '../models/progressMessage';

const tanzuDarwin = require('./tanzu-install-darwin.ts')
const tanzuWin32 = require('./tanzu-install-win32.ts')

function installUsingSteps(progressMessenger: ProgressMessenger) {
    let steps = []
    if (process.platform === 'darwin') {
        steps = tanzuDarwin.steps
    }
    return doInstall(progressMessenger, steps)
}

function doInstall(progressMessenger: ProgressMessenger, steps: InstallStep<ProgressMessenger>[]): boolean {
    // We cycle through the steps, executing each one (unless/until one returns FALSE) and we then return TRUE/FALSE
    return steps.reduce<boolean>((accumulatingResult, step) => {
        return accumulatingResult && step.execute(progressMessenger); // short-circuits if accumulatingResult is FALSE
    }, true )
}

if (process.platform === 'darwin') {
    module.exports = tanzuDarwin
} else if (process.platform === 'win32') {
    module.exports = tanzuWin32
}
module.exports.install = installUsingSteps

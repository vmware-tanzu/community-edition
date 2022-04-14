'use strict'

/**
 * The purpose of this module is to provide the main backend process (in electron.js) with OS-agnostic methods for doing all Concierge
 * activities. See the bottom of this module for methods that are exported.
 * Some methods are already OS-agnostic, so we simply export utility methods.
 * Other methods are OS-specific, so we choose our own implementation from our imports (using process.platform to decide which platform
 * we're on and therefore which implementation is appropriate)
 */


import { InstallationState, InstallData, InstallStep, PreInstallation } from '../models/installation';
import { ProgressMessenger } from '../models/progressMessage';

const tanzuDarwin = require('./tanzu-install-darwin.ts')
const tanzuWin32 = require('./tanzu-install-win32.ts')
const util = require('./tanzu-install-util.ts')

function installUsingSteps(preInstallation: PreInstallation, progressMessenger: ProgressMessenger) {
    return doInstall(preInstallation, progressMessenger, getInstallationData())
}

function doInstall(preInstallation: PreInstallation, progressMessenger: ProgressMessenger, installData: InstallData): InstallationState {
    util.reportInstallationStart(installData.msgStart, progressMessenger)
    // We cycle through the steps, executing each one (unless/until one returns an installation state with 'stop' set TRUE),
    // and we then return the final state from the last step
    const initialState = createInitialInstallationState(preInstallation, installData.steps);
    const finalState = installData.steps.reduce<InstallationState>(fxnExecuteStep(progressMessenger, installData.steps), initialState )
    const finalMessage = finalState.success ? installData.msgSucceeded : installData.msgFailed
    util.reportInstallationComplete(finalMessage, !finalState.success, progressMessenger)
    return finalState
}

// Returns a function that will execute a step (as part of doInstall's reduce() call above)
function fxnExecuteStep(progressMessenger: ProgressMessenger, steps: InstallStep[]) :
    (state: InstallationState, step: InstallStep, index: number) => InstallationState {
    return (previousInstallationState, step, index) => {
        // if the previous step returned an error and said STOP, we just return that state without executing this step
        if (previousInstallationState.stop) {
            return previousInstallationState;
        }
        // Create an installation state reflecting THIS step's data (and then execute the step)
        const installationStateParam = { ...previousInstallationState, currentStep: step.name, currentStepIndex: index };
        const newInstallationState = step.execute(installationStateParam, progressMessenger);

        // if we just successfully finished the last step in the steps array, we have succeeded
        const succeeded = !newInstallationState.stop && index === newInstallationState.totalSteps - 1;
        return { ...newInstallationState, success: succeeded };
    }
}

function createInitialInstallationState(preInstallation: PreInstallation, steps: InstallStep[]): InstallationState {
    const nSteps = steps ? steps.length : 0
    const existingInstallation = preInstallation.existingInstallation
    const chosenInstallation = preInstallation.chosenInstallation
    return { currentStep: firstStepName(steps), currentStepIndex: 0, totalSteps: nSteps, existingInstallation, chosenInstallation }
}

function getInstallationData() : InstallData {
    if (process.platform === 'darwin') {
        return tanzuDarwin.installData
    }
    return undefined
}

function firstStepName(steps: InstallStep[]) : string {
    return steps && steps.length > 0 ? steps[0].name : '';
}

if (process.platform === 'darwin') {
    module.exports.preinstall = tanzuDarwin.preinstall
} else if (process.platform === 'win32') {
    module.exports.preinstall = tanzuWin32.preinstall
}
module.exports.install = installUsingSteps
module.exports.pluginList = util.pluginList
module.exports.launchKickstart = util.launchKickstart
module.exports.launchTanzuUi = util.launchTanzuUi

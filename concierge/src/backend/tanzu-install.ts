'use strict'
import { ExistingInstallation, InstallationState, InstallStep } from '../models/installation';
import { ProgressMessenger } from '../models/progressMessage';

const tanzuDarwin = require('./tanzu-install-darwin.ts')
const tanzuWin32 = require('./tanzu-install-win32.ts')

function installUsingSteps(existingInstallation: ExistingInstallation, progressMessenger: ProgressMessenger) {
    return doInstall(existingInstallation, progressMessenger, getInstallationSteps())
}

function doInstall(existingInstallation: ExistingInstallation, progressMessenger: ProgressMessenger, steps: InstallStep[]): InstallationState {
    const initialState = createInitialInstallationState(existingInstallation, steps);
    // We cycle through the steps, executing each one (unless/until one returns an installation state with 'stop' set TRUE),
    // and we then return the final state from the last step
    return steps.reduce<InstallationState>(fxnExecuteStep(progressMessenger, steps), initialState );
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

function createInitialInstallationState(existingInstallation: ExistingInstallation, steps: InstallStep[]): InstallationState {
    const nSteps = steps ? steps.length : 0
    return { currentStep: firstStepName(steps), currentStepIndex: 0, totalSteps: nSteps, existingInstallation }
}

function getInstallationSteps() : InstallStep[] {
    if (process.platform === 'darwin') {
        return tanzuDarwin.steps
    }
    return undefined
}

function firstStepName(steps: InstallStep[]) : string {
    return steps && steps.length > 0 ? steps[0].name : '';
}

if (process.platform === 'darwin') {
    module.exports = tanzuDarwin
} else if (process.platform === 'win32') {
    module.exports = tanzuWin32
}
module.exports.install = installUsingSteps

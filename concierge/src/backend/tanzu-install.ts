'use strict'
import { ExistingInstallation, InstallationState, InstallStep } from '../models/installation';
import { ProgressMessenger } from '../models/progressMessage';

const tanzuDarwin = require('./tanzu-install-darwin.ts')
const tanzuWin32 = require('./tanzu-install-win32.ts')

function installUsingSteps(existingInstallation: ExistingInstallation, progressMessenger: ProgressMessenger) {
    return doInstall(existingInstallation, progressMessenger, getInstallationSteps())
}

function doInstall(existingInstallation: ExistingInstallation, progressMessenger: ProgressMessenger, steps: InstallStep[]): InstallationState {
    const initialState = createInitialInstallationState(existingInstallation, firstStepName(steps));
    // We cycle through the steps, executing each one (unless/until one returns an installation state with 'stop' set TRUE),
    // and we then return the final state from the last step
    return steps.reduce<InstallationState>(fxnExecuteStep(progressMessenger, steps), initialState );
}

// Returns a function that will execute a step (as part of the doInstall reduce() call)
function fxnExecuteStep(progressMessenger: ProgressMessenger, steps: InstallStep[]) :
    (state: InstallationState, step: InstallStep, index: number) => InstallationState {
    const nSteps = steps ? steps.length : 0;
    return (previousInstallationState, step, index) => {
        // if the previous step returned an error and said STOP, we just return that state without executing the step itself
        if (previousInstallationState.stop) {
            return previousInstallationState;
        }
        // Execute the step with THIS step's name
        const installationStateParam = { ...previousInstallationState, step: step.name };
        const newInstallationState = step.execute(installationStateParam, progressMessenger);

        // if we just successfully finished the last step in the steps array, we have succeeded
        const succeeded = !newInstallationState.stop && index === nSteps - 1;
        return { ...newInstallationState, success: succeeded };
    }
}

function createInitialInstallationState(existingInstallation: ExistingInstallation, firstStep: string): InstallationState {
    return { currentStep: firstStep, existingInstallation }
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

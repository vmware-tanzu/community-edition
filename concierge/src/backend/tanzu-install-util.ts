'use strict'
import { ProgressMessenger } from '../models/progressMessage';
import { InstallationState } from '../models/installation';

const { execSync } = require("child_process")
const fs = require('fs')
const yaml = require('js-yaml')

// Reports an installation-halting error and returns an installation state with stop===true
export function reportError(message: string, progressMessenger: ProgressMessenger, state: InstallationState) : InstallationState {
    progressMessenger.report({error: true, message, stepComplete: true, step: state.currentStep})
    return {...state, stop: true}
}
// Reports a high-level message of a current step
export function reportMessage(message: string, stepComplete: boolean, progressMessenger: ProgressMessenger, state: InstallationState) : InstallationState {
    progressMessenger.report({message, stepComplete, step: state.currentStep})
    return state
}
// Reports a step is complete (convenience method)
export function reportComplete(message: string, progressMessenger: ProgressMessenger, state: InstallationState) : InstallationState {
    return reportMessage(message, true, progressMessenger, state)
}
// Reports a low-level detail (message) of a current step
export function reportDetails(message = '', details: string, stepComplete: boolean, progressMessenger: ProgressMessenger, state: InstallationState) : InstallationState {
    progressMessenger.report({message, details, stepComplete, step: state.currentStep})
    return state
}
// Reports a percent complete, where percentComplete is 0 to 100
// Reports step is complete if percentComplete is 100
export function reportPercentComplete(percentComplete: number, progressMessenger: ProgressMessenger, state: InstallationState) {
    progressMessenger.report({message: '', percentComplete, stepComplete: percentComplete === 100, step: state.currentStep})
    return state
}

// Reports an installation-halting pre-requisite error and returns an installation state with stop===true
export function reportMissingPrerequisite(prerequisite: string, progressMessenger: ProgressMessenger, state: InstallationState) : InstallationState {
    const message = 'Some programmer forgot to ' + prerequisite + ' before starting the installation step "' + state.currentStep +
        '". We have to halt the installation. So sorry.'
    return reportError(message, progressMessenger, state)
}

// configPath should be:
// Darwin: os.homedir() + '/.config/tanzu/config.yaml'
export function tanzuEdition(configPath) {
    try {
        const fileContents = fs.readFileSync(configPath, 'utf8')
        const data = yaml.load(fileContents)
        if (!data) {
            console.log('opened config, but unable to parse data')
        } else {
            let edition = data?.clientOptions?.cli?.edition.toUpperCase()
            if (!edition) {
                console.log('Config file at ' + configPath + ' does not seem to include clientOptions.cli.edition')
            }
            let editionVersion = data?.clientOptions?.cli?.editionVersion
            if (!editionVersion) {
                console.log('Config file at ' + configPath + ' does not seem to include clientOptions.cli.editionVersion')
            }
            return {edition, editionVersion}
        }
    } catch (e) {
        console.log(e)
    }
    return {}
}

export function writeTanzuEdition(configPath, edition, version: string): boolean {
    try {
        const fileContents = fs.readFileSync(configPath, 'utf8')
        const data = yaml.load(fileContents)
        if (!data) {
            console.log('opened config, but unable to parse data')
        } else {
            if (!data?.clientOptions) {
                data.clientOptions = {}
            }
            if (!data.clientOptions.cli) {
                data.clientOptions.cli = {}
            }
            data.clientOptions.cli.edition = edition
            data.clientOptions.cli.editionVersion = version

            const newConfigData = yaml.dump(data)
            console.log('was thinking of writing: ' + newConfigData)
            fs.writeFileSync(configPath, newConfigData)
            return true
        }
    } catch (e) {
        console.log(e)
    }
    return false
}

// "command" should be "which" on Darwin/Linux and "where" for Windows
export function tanzuPath(command) {
    let path = ''
    const tanzuCommand = command + ' tanzu'
    try {
        const stdio = '' + execSync(tanzuCommand)
        const parts = stdio.match(/(.*)\n/)
        if (parts && parts.length > 1) {
            path = parts[1]
        } else {
            console.log('Unable to parse Tanzu path from output of "' + tanzuCommand + '" command: ' + stdio)
        }
    } catch(e) {
        console.log('Encountered error in executing ' + tanzuCommand  + ': ' + JSON.stringify(e))
    }
    return path
}

export function tanzuBinaryVersion() {
    let version = ''
    const cmd = 'tanzu version'
    try {
        const stdio = '' + execSync(cmd)
        const parts = stdio.match(/version: (.*)\n/)
        if (parts && parts.length > 1) {
            version = parts[1]
        } else {
            console.log('Unable to parse Tanzu version from output of ' + cmd + ' command: ' + stdio)
        }
    } catch(e) {
        console.log('Encountered error in executing ' + cmd  + ': ' + JSON.stringify(e))
    }
    return version
}

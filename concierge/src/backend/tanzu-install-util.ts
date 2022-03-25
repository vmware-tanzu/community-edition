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
// Reports a low-level detail (message) of a current step
export function reportDetails(message = '', details: string, stepComplete: boolean, progressMessenger: ProgressMessenger, state: InstallationState) : InstallationState {
    progressMessenger.report({message, details, stepComplete, step: state.currentStep})
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
        let fileContents = fs.readFileSync(configPath, 'utf8')
        const data = yaml.load(fileContents)
        if (!data) {
            console.log('opened config, but unable to parse data')
        } else {
            const edition = data?.clientOptions?.cli?.edition
            if (!edition) {
                console.log('Config file at ' + configPath + ' does not seem to include clientOptions.cli.edition: ' + JSON.stringify(data))
            } else {
                return edition.toUpperCase()
            }
        }
    } catch (e) {
        console.log(e)
    }
    return ''
}

// "command" should be "which" on Darwin/Linux and "where" for Windows
export function tanzuPath(command) {
    const tanzuCommand = command + ' tanzu'
    const stdio = '' + execSync(tanzuCommand)
    const parts = stdio.match(/(.*)\n/)
    let path = ''
    if (parts && parts.length > 1) {
        path = parts[1]
    } else {
        console.log('Unable to parse Tanzu path from output of "' + tanzuCommand + '" command: ' + stdio)
    }
    return path
}

export function tanzuVersion() {
    const stdio = '' + execSync('tanzu version')
    const parts = stdio.match(/version: (.*)\n/)
    let version = ''
    if (parts && parts.length > 1) {
        version = parts[1]
    } else {
        console.log('Unable to parse Tanzu version from output of "tanzu version" command: ' + stdio)
    }
    return version
}

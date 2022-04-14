'use strict'
import * as path from 'path';
import { ProgressMessage, ProgressMessenger } from '../models/progressMessage';
import { AvailableInstallation, InstallationArchive, InstallationState } from '../models/installation';
import { stringFromArray } from '_/utils';
const utils = require('../utils.ts')
const utilExec = require('./tanzu-exec-util.ts');

const { execSync } = require("child_process")
const fs = require('fs')
const yaml = require('js-yaml')

// Reports that the installation is starting
export function reportInstallationStart(message: string, progressMessenger: ProgressMessenger) {
    progressMessenger.report({message, installStarting: true})
}
// Reports that the installation has completed (error flag is set if failed)
export function reportInstallationComplete(message: string, error: boolean, progressMessenger: ProgressMessenger) {
    progressMessenger.report({message, installComplete: true, error})
}
// Reports an installation-halting error and returns an installation state with stop===true
export function reportError(message: string, progressMessenger: ProgressMessenger, state: InstallationState) : InstallationState {
    progressMessenger.report({error: true, message, stepComplete: true, step: state.currentStep})
    return {...state, stop: true}
}
// Reports a high-level message of a current step
export function reportMessage(message: string, progressMessenger: ProgressMessenger, state: InstallationState): InstallationState {
    progressMessenger.report({message, step: state.currentStep})
    return state
}
// Reports a step is complete (convenience method)
export function reportStepComplete(message: string, progressMessenger: ProgressMessenger, state: InstallationState) : InstallationState {
    return reportMessage(message, progressMessenger, state)
}
// Reports a step is starting (convenience method)
export function reportStepStart(message: string, progressMessenger: ProgressMessenger, state: InstallationState) : InstallationState {
    progressMessenger.report({message, stepStarting: true, step: state.currentStep})
    return state
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
export function tanzuEdition(configPath): {edition?: string, editionVersion?: string} {
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
export function tanzuPath(command) : ProgressMessage {
    let path = ''
    const tanzuCommand = command + ' tanzu'
    try {
        const stdio = '' + execSync(tanzuCommand)
        const parts = stdio.match(/(.*)\n/)
        if (parts && parts.length > 1) {
            path = parts[1]
        } else {
            // NOTE: this likely means there is no tanzu installed, rather than that we have an error
            const message = `Unable to parse Tanzu path from output of "${tanzuCommand}" command: ${stdio}`
            console.log(message)
            return { message }
        }
    } catch(e) {
        const message = `Encountered error in executing ${tanzuCommand}`
        const details = `Error was: ${JSON.stringify(e)}`
        console.log(`${message} ${details}`)
        return {message, details, error: true}
    }
    return { message: `path to tanzu is ${path}`, data: path }
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

// if the file exists, tries to remove the file and returns undefined if successful, an error message on failure
// if the file does not exist, just returns as if successful
export function removeFile(path: string): string {
    if (fs.existsSync(path)) {
        try {
            fs.unlinkSync(path)
        } catch (e) {
            const errMsg = `ERROR trying to remove file ${path}: ${e}`
            console.log(errMsg)
            return errMsg
        }
    }
}

export function pluginList(progressMessenger: ProgressMessenger): string[] {
    let result = []
    const pluginListResult = utilExec.exec({}, 'tanzu', 'plugin', 'list', '-o', 'json')
    if (pluginListResult.error) {
        progressMessenger.report(pluginListResult)
    } else {
        const pluginList = JSON.parse(pluginListResult.message)
        try {
            const plugins = pluginList as any[]
            return plugins.map<string>(p => p.name)
        } catch(e) {
            const message = 'Error trying to get list of plugins'
            progressMessenger.report({message, error: true, data: e})
        }
    }
    return result
}

export function launchKickstart(progressMessenger: ProgressMessenger) {
    const launchResult = utilExec.execAsync({stderrImpliesError: true}, 'tanzu', 'mc', 'create', '--ui')
    if (launchResult.error) {
        progressMessenger.report(launchResult)
    }
    progressMessenger.report({message: 'Kickstart UI launched', stepComplete: true})
}

export function launchTanzuUi(progressMessenger: ProgressMessenger) {
    const launchResult = utilExec.execAsync({stderrImpliesError: true}, 'tanzu', 'ui')
    if (launchResult.error) {
        progressMessenger.report(launchResult)
    }
    console.log('launchTanzuUiDarwin: ' + JSON.stringify(launchResult))
    progressMessenger.report({message: 'Tanzu UI launched', details: JSON.stringify(launchResult), stepComplete: true})
}

export function checkInstallationArchive(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    console.log('checkInstallationArchive...')
    reportStepStart('Checking for installation archive', progressMessenger, state)
    if (!state.chosenInstallation) {
        return reportMissingPrerequisite('set the chosen installation', progressMessenger, state)
    }

    const archiveDir = state.chosenInstallation.archive.dir
    const archiveFullPath = state.chosenInstallation.archive.fullPath

    if (!utils.pathExists(archiveDir)) {
        const message = 'ERROR: unfortunately, we are not able to find the directory where we expect an installation archive: ' + archiveDir + '' +
            ', so we\'ll have to abandon the installation effort. So sorry.'
        return reportError(message, progressMessenger, state)
    }
    progressMessenger.report({step: state.currentStep, message: 'Directory exists: ' + archiveDir})

    if (!utils.pathExists(archiveFullPath)) {
        const message = 'ERROR: unfortunately, we are not able to find the installation archive: ' + archiveFullPath + '' +
            ', so we\'ll have to abandon the installation effort. So sorry.'
        return reportError(message, progressMessenger, state)
    }

    return reportStepComplete('Found archive ' + archiveFullPath, progressMessenger, {...state})
}

export function unpackArchive(unpacker: (pathArchive, dirTmp: string) => ProgressMessage, state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    reportStepStart('Unpacking archive', progressMessenger, state)
    if (!state.chosenInstallation) {
        return reportMissingPrerequisite('set the chosen installation', progressMessenger, state)
    }
    if (!state.chosenInstallation.archive) {
        return reportMissingPrerequisite('detect archive information', progressMessenger, state)
    }
    if (!state.dirTanzuTmp) {
        return reportMissingPrerequisite('set the tanzu temp dir for unpacking the archive', progressMessenger, state)
    }
    const unpackResult = unpacker(state.chosenInstallation.archive, state.dirTanzuTmp)
    if (unpackResult.error) {
        const message = 'ERROR: unfortunately, we encountered an error trying to unpack the installation archive, ' +
            ' so we\'ll have to abandon the installation effort. So sorry.\n (Unpack from ' +
            state.chosenInstallation.archive.fullPath + ' to ' + state.dirTanzuTmp + ')'
        progressMessenger.report({step: state.currentStep, error: true, stepComplete: true,
            message, details: unpackResult.message + '\n' + unpackResult.details} )
        console.log('ERROR during unpack: ' + JSON.stringify(unpackResult))
        return {...state, stop: true}
    }
    console.log('SUCCESS unpack: ' + JSON.stringify(unpackResult))
    progressMessenger.report({...unpackResult, step: state.currentStep})

    const dirInstallFiles = path.join(state.dirTanzuTmp, expectedDirWithinArchive(state.chosenInstallation))
    // TODO: check that the expected dir actually exists
    const newState = {...state, dirInstallFiles}
    return reportStepComplete('Archive unpacked', progressMessenger, newState)
}

function expectedDirWithinArchive(installation: AvailableInstallation) {
    // return `${installation.edition}-${installation.machineArchitecture}-${installation.version}`
    return installation.mainDirAfterUnpack
}

// returns a regular expression to match an archive and pull out the edition and version
export function regExpMatchArchive(machineArchitecture, extension: string): RegExp {
    return new RegExp(`^([^-]*)-${machineArchitecture}-(v[\\d\\.]+)\\.${extension}$`)
}

// given a root directory and an array of archive files found IN that directory, create an array of AvailableInstallation objects.
// To do that, we need a regular expression to parse the archive name into edition+version, and we need a way of getting the main dir
// that the archive will expand (where we would find the tanzu binary and plugins)
export function createAvailableInstallations(dir:string, archives: string[], regexp: RegExp,
                                           mainArchiveDir: (edition: string, file: string) => string ): AvailableInstallation[] {
    const result = archives.map<AvailableInstallation>(archive => {
        const arrayTarballParts = archive.match(regexp)
        const file = stringFromArray(arrayTarballParts, 0)
        const edition = stringFromArray(arrayTarballParts, 1)
        const version = stringFromArray(arrayTarballParts, 2)
        const mainDirAfterUnpack = mainArchiveDir(edition, version)
        return {version, archive: {dir, file, fullPath: path.join(dir, file) }, edition, mainDirAfterUnpack }
    } )
    return result
}

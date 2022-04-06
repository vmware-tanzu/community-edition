import {
    AvailableInstallation,
    ExistingInstallation,
    InstallationState,
    InstallationArchive, InstallData,
    PreInstallation
} from '../models/installation'
import { ProgressMessage, ProgressMessenger } from '../models/progressMessage'
import * as path from 'path';

const { spawn, spawnSync } = require("child_process")
const os = require( 'os' )
const fs = require('fs')

const utilInstall = require('./tanzu-install-util.ts');
const utilExec = require('./tanzu-exec-util.ts');
const utils = require('../utils.ts')

const darwinSteps = [
    { name: 'Check prerequisites', execute: utilInstall.checkInstallationArchive },
    { name: 'Set data paths', execute: darwinSetDataPaths },
    { name: 'Unpack tanzu', execute: darwinTarballUnpack },
    { name: 'Set binary path', execute: darwinSetBinPath },
    { name: 'Delete existing Tanzu binary', execute: darwinDeleteExistingTanzuIfNec },
    { name: 'Copy new Tanzu binary', execute: darwinCopyTanzuBinary },
    { name: 'Copy uninstall script', execute: darwinCopyUninstallScript },
    { name: 'Copy plugins', execute: darwinCopyPlugins },
    { name: 'Install plugins', execute: darwinInstallPlugins },
    { name: 'Add repos', execute: darwinAddTanzuReposIfNec },
    // NOTE: we set the edition in the config file AFTER installing the plugins so that the tanzu binary will have created the config file
    // TODO: only set the edition in the config file if the tanzu installation supports the command
    { name: 'Set edition in config file', execute: darwinSetEdition },
]

function preinstallDarwin(progressMessenger: ProgressMessenger): PreInstallation {
    const fixPathResult = utils.fixPath()
    progressMessenger.report({message: fixPathResult.message, details: fixPathResult.techMessage, warning: fixPathResult.error})

    const existingInstallation = detectExistingInstallation(darwinConfigPath(), progressMessenger)
    const dirInstallationTarballsExpected = getInstallationDirs()

    dirInstallationTarballsExpected.forEach(dir => {
        progressMessenger.report({step: 'PRE-INSTALL', message: `files in dir ${dir}: [${listFiles(dir).join('][')}]`})
    })

    const availableInstallations = detectAvailableInstallations(dirInstallationTarballsExpected)
    return { existingInstallation, availableInstallations, dirInstallationTarballsExpected }
}

function getInstallationDirs(): string[] {
    const localDir = path.join(process.resourcesPath, '..', 'tanzu-releases')
    const userDir = path.join(os.homedir(), 'tanzu-releases')
    return [localDir, userDir]
}

function detectAvailableInstallations(dirs: string[]): AvailableInstallation[] {
    return dirs.reduce<AvailableInstallation[]>((accum, dir) => {
        accum.push(...detectAvailableInstallationsInDir(dir));
        return accum;
    }, []);
}

function detectAvailableInstallationsInDir(dir: string): AvailableInstallation[] {
    const machineArchitecture = 'darwin-amd64'
    // NOTE: we're looking for files with a name like: tce-darwin-amd64-v0.11.0.tar.gz where edition=tce and version=v0.11.0
    console.log(`Detecting available installation by looking in dir ${dir} for tarballs`)
    // TODO: use machineArchitecture in RegEx and move to util
    const tarballs = listFilesFiltered(dir, /^[^-]*-darwin-amd64-v[\d\.]+\.tar\.gz$/)
    console.log(`TARBALLS: [${tarballs.join('], [')}]`)
    const result = tarballs.map<AvailableInstallation>(tarball => {
        // this should always match, due to expression above
        const arrayTarballParts = tarball.match(/^([^-]*)-darwin-amd64-(v[\d\.]+)\.tar\.gz$/)
        const file = arrayTarballParts ? arrayTarballParts[0] : ''
        const edition = arrayTarballParts ? arrayTarballParts[1] : ''
        const version = arrayTarballParts ? arrayTarballParts[2] : ''
        return {version, archive: {dir, file, fullPath: dir + '/' + file }, edition, machineArchitecture }
    } )
    return result
}

function detectExistingInstallation(configPath: string, progressMessenger: ProgressMessenger): ExistingInstallation {
    const pathResult = utilInstall.tanzuPath('which')
    progressMessenger.report(pathResult)
    if (pathResult.data) {
        const tanzuBinaryVersion = utilInstall.tanzuBinaryVersion()
        const editionResult = utilInstall.tanzuEdition(configPath)
        const result = {path: pathResult.data, tanzuBinaryVersion, edition: editionResult.edition, editionVersion: editionResult.editionVersion}
        console.log('Existing install: ' + JSON.stringify(result))
        return result
    }
    return undefined
}

function darwinTarballUnpack(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    return utilInstall.unpackArchive(darwinUntar, state, progressMessenger)
}

function darwinSetDataPaths(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    const dirTanzuData = os.homedir() + '/.local/share'
    const pathTanzuConfig = darwinConfigPath()
    const dirTanzuConfig = darwinConfigDir()
    const dirTanzuTmp = os.homedir() + '/.tanzu/concierge/tmp'

    const tmpResult = darwinEnsureDir(dirTanzuTmp)
    if (tmpResult.error) {
        const message = `Unable to ensure tmp dir (for untar-ing) ${dirTanzuTmp}`
        console.log(`ERROR: ${message}; raw command result ${JSON.stringify(tmpResult)}`)
        return utilInstall.reportError(message, progressMessenger, state)
    }

    const newState = {...state, pathTanzuConfig, dirTanzuData, dirTanzuConfig, dirTanzuTmp}
    return utilInstall.reportStepComplete('Set data paths', progressMessenger, newState)
}

function darwinSetBinPath(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    const resultEnvPath = darwinExec('echo', '$PATH')
    if (resultEnvPath.error) {
        progressMessenger.report(resultEnvPath)
        return {...state, stop: true}
    }
    const envPath = resultEnvPath.message
    const preferredPath = os.homedir() + '/bin:'
    const defaultPath = '/usr/local/bin'
    const dirTanzuInstallation = darwinEnvPathContains(envPath, preferredPath) ? preferredPath : defaultPath
    const newState = {...state, dirTanzuInstallation}
    return utilInstall.reportDetails('Set binary path', 'Binary path=' + dirTanzuInstallation, true, progressMessenger, newState)
}

function darwinDeleteExistingTanzuIfNec(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.existingInstallation) {
        return utilInstall.reportStepComplete('No existing tanzu binary found (so skipping delete)', progressMessenger, state)
    }
    const result = darwinExec('rm', '-f', state.existingInstallation.path)
    if (result.error) {
        return utilInstall.reportError(result.message, progressMessenger, state)
    }
    const message =  'Successful removal of ' + state.existingInstallation.path
    return utilInstall.reportStepComplete(message, progressMessenger, state)
}

function darwinCopyTanzuBinary(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    utilInstall.reportStepStart('Copy tanzu binary', progressMessenger, state)
    if (!state.dirInstallFiles) {
        return utilInstall.reportMissingPrerequisite('set the source binary directory', progressMessenger, state)
    }
    if (!state.dirTanzuInstallation) {
        return utilInstall.reportMissingPrerequisite('set the target binary directory', progressMessenger, state)
    }

    const pathSourceBinary = state.dirInstallFiles + '/tanzu'
    if (!fs.existsSync(pathSourceBinary)) {
        const msg = 'Tanzu binary (' + pathSourceBinary +
            ') does not exist. (Is it possible the tarball was incomplete, malformed or in a new format?)'
        return utilInstall.reportError(msg, progressMessenger, state)
    }

    const pathTargetBinary = state.dirTanzuInstallation + '/tanzu'
    const result = darwinExec('install', pathSourceBinary, pathTargetBinary)
    if (result.error) {
        return utilInstall.reportError(result.message, progressMessenger, state)
    }
    const message =  `Successful copy of tanzu binary (${pathSourceBinary}) to ${pathTargetBinary}`
    return utilInstall.reportStepComplete(message, progressMessenger, state)
}

function darwinSetEdition(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    utilInstall.reportStepStart('Setting edition', progressMessenger, state)
    if (!state.chosenInstallation) {
        return utilInstall.reportMissingPrerequisite('set the chosen installation', progressMessenger, state)
    }
    if (!state.chosenInstallation.edition) {
        return utilInstall.reportMissingPrerequisite('set edition of the chosen installation', progressMessenger, state)
    }
    const edition = state.chosenInstallation.edition
    const version = state.chosenInstallation.version

    let message = ''
    if (utilInstall.writeTanzuEdition(state.pathTanzuConfig, edition, version)) {
        message = 'Set edition and version in config file (to ' + edition.toUpperCase() + ' ' + version + ')'
    } else {
        message = 'Unable to update config with edition and version. You should still be able to run Tanzu, but it\'s disappointing.'
    }
    return utilInstall.reportStepComplete(message, progressMessenger, state)
}

function darwinCopyUninstallScript(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    utilInstall.reportStepStart('Copy uninstall script', progressMessenger, state)
    if (!state.dirInstallFiles) {
        return utilInstall.reportMissingPrerequisite('set the source binary directory', progressMessenger, state)
    }
    if (!state.dirTanzuData) {
        return utilInstall.reportMissingPrerequisite('set the tanzu data directory', progressMessenger, state)
    }
    const tceDataDir = state.dirTanzuData + '/tce'
    const resultMakeTceDir = darwinExec('mkdir', '-p', tceDataDir)
    if (resultMakeTceDir.error) {
        return utilInstall.reportError(resultMakeTceDir.message, progressMessenger, state)
    }
    // TODO: log detail of creating/ensuring tce dir

    const sourceScript = state.dirInstallFiles + '/uninstall.sh'
    const resultInstallScript = darwinExec('install', sourceScript, tceDataDir)
    if (resultInstallScript.error) {
        return utilInstall.reportError(resultInstallScript.message, progressMessenger, state)
    }
    const message = 'copied uninstall script to ' + tceDataDir
    return utilInstall.reportStepComplete(message, progressMessenger, state)
}

// copies the plugins (and discovery info) from tarball dir to tanzu config dir
// returns state with list of plugins added
function darwinCopyPlugins(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    utilInstall.reportStepStart('Copy plugins', progressMessenger, state)
    // remove the old plugin cache so new plugins are detected
    const pluginCacheFile = os.homedir() + '/.cache/tanzu/catalog.yaml'
    const errMsg = utilInstall.removeFile(pluginCacheFile)
    if (errMsg) {
        return utilInstall.reportError(errMsg, progressMessenger, state)
    }
    utilInstall.reportMessage(`Deleted ${pluginCacheFile}`, progressMessenger, state)

    const dirDstPlugin = state.dirTanzuConfig + '/tanzu-plugins'
    // copy the discovery dir
    const dirSrcPluginDiscovery = state.dirInstallFiles + '/default-local/discovery'
    const resultDiscoveryCopy = darwinCopyRecursive(dirSrcPluginDiscovery, dirDstPlugin)
    if (resultDiscoveryCopy.error) {
        return utilInstall.reportError(resultDiscoveryCopy.message, progressMessenger, state)
    }
    utilInstall.reportMessage(`Copied discovery directory ${dirSrcPluginDiscovery}`, progressMessenger, state)

    // copy the plugins themselves
    const dirSrcPluginDistribution = state.dirInstallFiles + '/default-local/distribution/darwin/amd64/cli'
    const dirDstPluginDistribution = dirDstPlugin + '/distribution/darwin/amd64/cli'
    const plugins = listFiles(dirSrcPluginDistribution)
    const nPlugins = plugins.length
    let errorState
    plugins.every( (plugin, index) => {
        const dirSrcThisPlugin = dirSrcPluginDistribution + '/' + plugin
        const resultPluginCopy = darwinCopyRecursive(dirSrcThisPlugin, dirDstPluginDistribution)
        if (resultPluginCopy.error) {
            errorState = utilInstall.reportError(resultPluginCopy.message, progressMessenger, state)
            return false
        } else {
            console.log(`PLUGIN COPY RESULT: ${JSON.stringify(resultPluginCopy)}`)
        }
        utilInstall.reportDetails('', 'Copied plugin ' + plugin, false, progressMessenger, state)
        const percentComplete = utils.percentage(index + 1, nPlugins)
        utilInstall.reportPercentComplete(percentComplete, progressMessenger, state)
        return true
    })
    if (errorState) {
        return errorState
    }

    const newState = {...state, plugins}
    const msgSuccess = `Completed copying ${nPlugins} plugins`
    return utilInstall.reportStepComplete(msgSuccess, progressMessenger, newState)
}

function darwinInstallPlugins(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    const nPlugins = state.plugins?.length
    utilInstall.reportStepStart(`Install ${nPlugins} plugins`, progressMessenger, state)

    if (!state.plugins) {
        return utilInstall.reportMissingPrerequisite('set the plugin list', progressMessenger, state)
    }

    let errorState
    state.plugins.every( (plugin, index) => {
        const resultPluginInstall = darwinPluginInstall(plugin)
        if (resultPluginInstall.error) {
            errorState = utilInstall.reportError(resultPluginInstall.message, progressMessenger, state)
            return false
        } else {
            console.log(`PLUGIN INSTALL RESULT for ${plugin}: ${JSON.stringify(resultPluginInstall)}`)
        }
        utilInstall.reportDetails('', 'Installed plugin ' + plugin, false, progressMessenger, state)
        utilInstall.reportPercentComplete(utils.percentage(index + 1, nPlugins), progressMessenger, state)
        return true
    })
    if (errorState) {
        return errorState
    }

    return utilInstall.reportStepComplete('Completed plugin installation', progressMessenger, state)
}

function darwinAddTanzuReposIfNec(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    utilInstall.reportStepStart('Adding Tanzu plugin repos', progressMessenger, state)
    progressMessenger.report(darwinEnsureTanzuRepo('tce', 'tce-tanzu-cli-plugins', 'artifacts'))
    progressMessenger.report(darwinEnsureTanzuRepo('core-admin', 'tce-tanzu-cli-framework-admin', 'artifacts-admin'))
    return utilInstall.reportStepComplete('Done adding Tanzu repos for TCE', progressMessenger, state)
}

//============================================================
// helper methods
//============================================================
function darwinConfigPath(): string {
    return darwinConfigDir() + '/tanzu/config.yaml'
}

function darwinConfigDir(): string {
    return os.homedir() + '/.config'
}

function listFiles(srcDir: string): string[] {
    const resultListing = darwinExec('ls', '-a', srcDir)
    // filter out '.' and '..' and ending blank line
    const files = resultListing.message.split('\n').filter(listing => listing.length > 2)
    console.log('files in dir ' + srcDir + ' yields [' + files.join('],[') + ']')
    return files
}

function listFilesFiltered(srcDir: string, regex) {
    return listFiles(srcDir).filter(file => file.match(regex) !== null)
}

function darwinPluginInstall(plugin: string): ProgressMessage {
    const installResult = darwinExec('tanzu', 'plugin', 'install', plugin)
    // workaround: we detect a tanzu error in the error message to stderr, even if the exit code didn't indicate an error
    if (!installResult.error && installResult.details && installResult.details.startsWith('Error: ')) {
        return { ...installResult, error: true }
    }
    return installResult
}

function darwinEnsureTanzuRepo(repo, gcpBucket, gcpRootpath): ProgressMessage {
    if (darwinTanzuRepoExists(repo)) {
        return {message: `Repo ${repo} already exists`}
    }
    return darwinExec('tanzu', 'plugin', 'repo', 'add', '--name', repo,
        '--gcp-bucket-name', gcpBucket, '--gcp-root-path', gcpRootpath)
}

function darwinTanzuRepoExists(repo): boolean {
    const repoListReturn = darwinExec('tanzu', 'plugin', 'repo', 'list')
    console.log(`Result of repo listing: ${JSON.stringify(repoListReturn)}`)
    return !repoListReturn.error && repoListReturn.message && repoListReturn.message.includes(repo)
}

function darwinExec(command: string, ...args: string[]) : ProgressMessage {
    return doDarwinExec(spawnSync, command, ...args)
}

function doDarwinExec(fxn: any, command: string, ...args: string[]) : ProgressMessage {
    console.log(`doDarwinExec(): ${command} ${args.join(' ')}`)
    const result = {message: '', details: '', error: false }
    try {
        const syncResult = fxn(command, args, {stdio: 'pipe', encoding: 'utf8'})
        console.log(`doDarwinExec(${command}) yields: ${JSON.stringify(syncResult)}`)
        result.message = syncResult.stdout?.toString()
        result.details = syncResult.stderr?.toString()
    } catch (e) {
        console.log(e)
        result.message = 'ERROR: ' + e.toString()
        result.error = true
    }
    return result
}

function darwinEnsureDir(dir: string): ProgressMessage {
    return utilExec.exec({},'mkdir', '-p', dir)
}

// convenience wrapper to darwinExec for a copy command
function darwinCopyRecursive(src, dst: string): ProgressMessage {
    console.log(`Ensuring dir ${dst} exists`)
    const resultEnsureDir = darwinEnsureDir(dst)
    if (resultEnsureDir.error) {
        return resultEnsureDir
    }
    console.log(`Trying recursive copy from ${src} to ${dst}`)
    return darwinExec('cp', '-r', src, dst)
}

function darwinUntar(tarballInfo: InstallationArchive, dstDir: string) : ProgressMessage {
    const result = utilExec.exec({}, 'tar', 'xzvf', tarballInfo.fullPath, '-C', dstDir)
    if (!result.error) {
        result.message =  'Successful untar of ' + tarballInfo.fullPath
    }
    result.stepComplete = true
    return result
}

function darwinEnvPathContains(envPath, target) {
    return utils.stringContains(':' + envPath + ':', ':' + target + ':')
}

function expectedDirWithinTarball(installation: AvailableInstallation) {
    return `${installation.edition}-darwin-amd64-${installation.version}`
}

const installData = {
    steps: darwinSteps,
    msgStart: 'Here we go... (starting installation on Mac OS)',
    msgFailed: 'So sorry the installation did not succeed. Please try again after any known issues are addressed.',
    msgSucceeded: 'You\'re ready to start using Tanzu!',
} as InstallData
module.exports.installData = installData
module.exports.preinstall = preinstallDarwin

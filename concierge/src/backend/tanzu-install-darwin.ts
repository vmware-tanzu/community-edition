import {
    AvailableInstallation,
    ExistingInstallation,
    InstallationState,
    InstallationTarball,
    PreInstallation
} from '../models/installation'
import { ProgressMessage, ProgressMessenger } from '../models/progressMessage'

const { spawnSync } = require("child_process")
const os = require( 'os' )
const fs = require('fs')

const util = require('./tanzu-install-util.ts');
const utils = require('../utils.ts')

const darwinSteps = [
    { name: 'Check prerequisites', execute: darwinTarballCheck },
    { name: 'Unpack tanzu', execute: darwinTarballUnpack },
    { name: 'Set data paths', execute: darwinSetDataPaths },
    { name: 'Set binary path', execute: darwinSetBinPath },
    { name: 'Delete existing Tanzu binary', execute: darwinDeleteExistingTanzuIfNec },
    { name: 'Copy new Tanzu binary', execute: darwinCopyTanzuBinary },
    // TODO: only set the edition in the config file if the tanzu installation supports the command
    { name: 'Set edition in config file', execute: darwinSetEdition },
    { name: 'Copy uninstall script', execute: darwinCopyUninstallScript },
    { name: 'Copy plugins', execute: darwinCopyPlugins },
    { name: 'Install plugins', execute: darwinInstallPlugins },
    { name: 'Add repos', execute: darwinAddTanzuReposIfNec },
]

function darwinConfigPath(): string {
    return darwinConfigDir() + '/config.yaml'
}

function darwinConfigDir(): string {
    return os.homedir() + '/.config/tanzu'
}

function preinstallDarwin(): PreInstallation {
    const existingInstallation = detectExistingInstallation(darwinConfigPath())
    const availableInstallations = detectAvailableInstallations()

    return { existingInstallation, availableInstallations }
}

function detectAvailableInstallations(): AvailableInstallation[] {
    // NOTE: we're looking for files with a name like: tce-darwin-amd64-v0.11.0.tar.gz where edition=tce and version=v0.11.0
    const dir = __dirname
    const tarballs = listFilesFiltered(dir, /^tce-darwin-amd64-v[\d\.]+\.tar\.gz$/)
    console.log(`TARBALLS: [${tarballs.join('], [')}]`)
    const result = tarballs.map<AvailableInstallation>(tarball => {
        const arrayTarballParts = tarball.match(/^([^-]*)-darwin-amd64-(v[\d\.]+)\.tar\.gz$/)
        const file = arrayTarballParts[0]
        const edition = arrayTarballParts[1]
        const version = arrayTarballParts[2]
        return {version, tarball: {dir, file, fullPath: dir + '/' + file }, edition}
    } )
    return result
}

function detectExistingInstallation(configPath: string): ExistingInstallation {
    const path = util.tanzuPath('which')
    if (path) {
        const tanzuBinaryVersion = util.tanzuBinaryVersion()
        const editionResult = util.tanzuEdition(configPath)
        const result = {path, tanzuBinaryVersion, edition: editionResult.edition, editionVersion: editionResult.editionVersion}
        console.log('Existing install: ' + JSON.stringify(result))
        return result
    }
    return undefined
}

function darwinTarballCheck(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    console.log('darwinTarballCheck...')
    // TODO: move this "here we go..." to another step?
    progressMessenger.report({message: 'Here we go... (starting installation on Mac OS)'})

    if (!state.chosenInstallation) {
        return util.reportMissingPrerequisite('set the chosen installation', progressMessenger, state)
    }

    const tarballDir = state.chosenInstallation.tarball.dir
    // const tarballFile = state.chosenInstallation.tarball.file
    const tarballFullPath = state.chosenInstallation.tarball.fullPath

    if (!darwinTestPath(tarballDir)) {
        const message = 'ERROR: unfortunately, we are not able to find the directory where we expect an installation tarball: ' + tarballDir + '' +
        ', so we\'ll have to abandon the installation effort. So sorry.'
        return util.reportError(message, progressMessenger, state)
    }
    progressMessenger.report({step: state.currentStep, message: 'Directory exists: ' + tarballDir})

    if (!darwinTestPath(tarballFullPath)) {
        const message = 'ERROR: unfortunately, we are not able to find the installation tarball: ' + tarballFullPath + '' +
            ', so we\'ll have to abandon the installation effort. So sorry.'
        return util.reportError(message, progressMessenger, state)
    }

    return util.reportMessage('Found tarball ' + tarballFullPath, true, progressMessenger, {...state})
}

function darwinTarballUnpack(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.chosenInstallation) {
        return util.reportMissingPrerequisite('set the chosen installation', progressMessenger, state)
    }
    if (!state.chosenInstallation.tarball) {
        return util.reportMissingPrerequisite('detect tarball information', progressMessenger, state)
    }
    const untarResult = darwinUntar(state.chosenInstallation.tarball)
    if (untarResult.error) {
        progressMessenger.report({step: state.currentStep, error: true, stepComplete: true,
            message: 'ERROR: unfortunately, we encountered an error trying to untar the installation tarball: ' +
                state.chosenInstallation.tarball.fullPath + ', so we\'ll have to abandon the installation effort. So sorry.',
            details: untarResult.message + '\n' + untarResult.details} )
        console.log('ERROR during untar: ' + JSON.stringify(untarResult))
        return {...state, stop: true}
    }
    console.log('SUCCESS untar: ' + JSON.stringify(untarResult))
    progressMessenger.report({...untarResult, step: state.currentStep})

    // TODO: dynamically find the untarred path
    //const hardCodedPath = '/Users/swalner/workspace/tanzu-community-edition/community-edition/concierge/dist/tce-darwin-amd64-v0.10.0'
    const hardCodedPath = '/Users/swalner/workspace/tanzu-community-edition/community-edition/concierge/dist/tce-darwin-amd64-v0.11.0'
    return {...state, dirInstallFiles: hardCodedPath}
}

function darwinSetDataPaths(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    const dirTanzuData = os.homedir() + '/.local/share'
    const pathTanzuConfig = darwinConfigPath()
    const dirTanzuConfig = darwinConfigDir()

    const newState = {...state, pathTanzuConfig, dirTanzuData, dirTanzuConfig}
    return util.reportMessage('Set data paths', true, progressMessenger, newState)
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
    return util.reportDetails('Set binary path', 'Binary path=' + dirTanzuInstallation, true, progressMessenger, newState)
}

function darwinDeleteExistingTanzuIfNec(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.existingInstallation) {
        return util.reportMessage('No existing tanzu binary found (so skipping delete)', true, progressMessenger, state)
    }
    const result = darwinExec('rm', '-f', state.existingInstallation.path)
    if (result.error) {
        return util.reportError(result.message, progressMessenger, state)
    }
    const message =  'Successful removal of ' + state.existingInstallation.path
    return util.reportMessage(message, true, progressMessenger, state)
}

function darwinCopyTanzuBinary(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.dirInstallFiles) {
        return util.reportMissingPrerequisite('set the source binary directory', progressMessenger, state)
    }
    if (!state.dirTanzuInstallation) {
        return util.reportMissingPrerequisite('set the target binary directory', progressMessenger, state)
    }

    const pathSourceBinary = state.dirInstallFiles + '/tanzu'
    if (!fs.existsSync(pathSourceBinary)) {
        const msg = 'Tanzu binary (' + pathSourceBinary +
            ') does not exist. (Is it possible the tarball was incomplete, malformed or in a new format?)'
        return util.reportError(msg, progressMessenger, state)
    }

    const pathTargetBinary = state.dirTanzuInstallation + '/tanzu'
    const result = darwinExec('install', pathSourceBinary, pathTargetBinary)
    if (result.error) {
        return util.reportError(result.message, progressMessenger, state)
    }
    const message =  'Successful copy of tanzu binary to ' + pathTargetBinary
    const detail = 'Source binary: ' + pathSourceBinary
    return util.reportDetails(message, detail, true, progressMessenger, state)
}

function darwinSetEdition(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.chosenInstallation) {
        return util.reportMissingPrerequisite('set the chosen installation', progressMessenger, state)
    }
    if (!state.chosenInstallation.edition) {
        return util.reportMissingPrerequisite('set edition of the chosen installation', progressMessenger, state)
    }
    const edition = state.chosenInstallation.edition
    const version = state.chosenInstallation.version
/*
    const resultSetEdition = darwinExec('tanzu', ['config', 'set', 'env.edition', edition])
    if (resultSetEdition.error) {
        return util.reportError(resultSetEdition.message, progressMessenger, state)
    }
*/

    let message = ''
    if (util.writeTanzuEdition(state.pathTanzuConfig, edition, version)) {
        message = 'Set edition and version in config file (to ' + edition.toUpperCase() + ' ' + version + ')'
    } else {
        message = 'Unable to update config with edition and version. You should still be able to run Tanzu, but it\'s disappointing.'
    }
    return util.reportMessage(message, true, progressMessenger, state)
}

function darwinCopyUninstallScript(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    if (!state.dirInstallFiles) {
        return util.reportMissingPrerequisite('set the source binary directory', progressMessenger, state)
    }
    if (!state.dirTanzuData) {
        return util.reportMissingPrerequisite('set the tanzu data directory', progressMessenger, state)
    }
    const tceDataDir = state.dirTanzuData + '/tce'
    const resultMakeTceDir = darwinExec('mkdir', '-p', tceDataDir)
    if (resultMakeTceDir.error) {
        return util.reportError(resultMakeTceDir.message, progressMessenger, state)
    }
    // TODO: log detail of creating/ensuring tce dir

    const sourceScript = state.dirInstallFiles + '/uninstall.sh'
    const resultInstallScript = darwinExec('install', sourceScript, tceDataDir)
    if (resultInstallScript.error) {
        return util.reportError(resultInstallScript.message, progressMessenger, state)
    }
    const message = 'copied uninstall script to ' + tceDataDir
    return util.reportComplete(message, progressMessenger, state)
}

// copies the plugins (and discovery info) from tarball dir to tanzu config dir
// returns state with list of plugins added
function darwinCopyPlugins(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    const dirSrcPluginDiscovery = state.dirInstallFiles + '/default-local/discovery'
    const dirDstPlugin = state.dirTanzuConfig + '/tanzu-plugins'

    // copy the discovery dir
    const resultDiscoveryCopy = darwinCopyRecursive(dirSrcPluginDiscovery, dirDstPlugin)
    if (resultDiscoveryCopy.error) {
        return util.reportError(resultDiscoveryCopy.message, progressMessenger, state)
    }
    util.reportMessage('Finished copying discovery dir', false, progressMessenger, state)

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
            errorState = util.reportError(resultPluginCopy.message, progressMessenger, state)
            return false
        } else {
            console.log(`PLUGIN COPY RESULT: ${JSON.stringify(resultPluginCopy)}`)
        }
        util.reportDetails('', 'Copied plugin ' + plugin, false, progressMessenger, state)
        const percentComplete = utils.percentage(index + 1, nPlugins)
        util.reportPercentComplete(percentComplete, progressMessenger, state)
        return true
    })
    if (errorState) {
        return errorState
    }

    const newState = {...state, plugins}
    const msgSuccess = `Completed copying ${nPlugins} plugins`
    return util.reportComplete(msgSuccess, progressMessenger, newState)
}

function darwinInstallPlugins(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    // install the plugins
    const message = 'PRETENDED TO INSTALL PLUGINS: ' + state.plugins.join(',')
    return util.reportComplete(message, progressMessenger, state)
}

function XdarwinInstallPlugins(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    const pluginSrcDir = state.dirInstallFiles + '/bin'
    const plugins = listFiles(pluginSrcDir)
    const nPlugins = plugins.length
    const pluginDstDir = state.dirTanzuConfig + '/tanzu-plugins'

    util.reportMessage(`Preparing to install ${nPlugins} plugins` , false, progressMessenger, state)

    let errorState
    plugins.every( (plugin, index) => {
        const resultPluginInstall = tanzuPluginInstall(pluginSrcDir, plugin, pluginDstDir)
        if (resultPluginInstall.error) {
            errorState = util.reportError(resultPluginInstall.message, progressMessenger, state)
            return false
        } else {
            console.log(`PLUGIN INSTALL RESULT: ${JSON.stringify(resultPluginInstall)}`)
        }
        util.reportDetails('', 'Installed plugin ' + plugin, false, progressMessenger, state)
        const percentComplete = utils.percentage(index + 1, nPlugins)
        util.reportPercentComplete(percentComplete, progressMessenger, state)
        return true
    })
    if (errorState) {
        return errorState
    }

    return util.reportMessage('Completed plugin installation', true, progressMessenger, state)
}

function darwinAddTanzuReposIfNec(state: InstallationState, progressMessenger: ProgressMessenger) : InstallationState {
    // TODO: implement adding TCE repos
    const message = 'PRETENDED TO ADD TANZU REPO'
    return util.reportMessage(message, true, progressMessenger, state)
}

//============================================================
// helper methods
//============================================================
function listFiles(srcDir: string): string[] {
    const resultListing = darwinExec('ls', '-a', srcDir)
    // filter out '.' and '..' and ending blank line
    const files = resultListing.message.split('\n').filter(listing => listing.length > 2)
    console.log('files in dir ' + srcDir + ' yields [' + files.join('],[') + ']')
    return files
}

function listFilesFiltered(srcDir: string, regex) {
    return listFiles(srcDir).filter(file => {
        const regexArray = file.match(regex)
        const result = regexArray !== null
        console.log(`FILE ${file} gets result ${result} from regex ${regex}`)
        return result
    })
}

// srcDir is the directory in which the plugin is found (where the tarball was expanded + /bin)
// srcName is the name of the binary, expected to be tanzu-plugin-xxx
// dstDir is the plugin directory under the config area, generally CONFIG/tanzu-plugins
export function tanzuPluginInstall(srcDir, srcName, dstDir): ProgressMessage {
    // COPY the src plugin to the dstDir
    const resultCopy = darwinCopyPlugin(srcDir, srcName, dstDir)
    console.log(`darwinCopyPlugin just returned ${JSON.stringify(resultCopy)}`)
    if (resultCopy.error) {
        return resultCopy
    }

    const pluginName = parsePluginName(srcName)
    if (!pluginName) {
        return { error: true, message: 'Unable to parse plugin name ' + pluginName}
    } else {
        console.log(`PLUGIN NAME: ${pluginName}`)
    }

    return darwinExec('tanzu', 'plugin', 'install', pluginName)
}

function darwinCopyPlugin(srcDir, srcName, dstDir: string): ProgressMessage {
    const pathSrcPlugin = srcDir + '/' + srcName
    console.log(`About to copy from ${pathSrcPlugin} to ${dstDir}`)
    return darwinExec('cp', pathSrcPlugin, dstDir)
}

// Take a plugin filename like tanzu-plugin-xxx and return the plugin name xxx
function parsePluginName(pluginName: string): string {
    const START_OF_NAME = 'tanzu-plugin-'
    if (!pluginName.startsWith(START_OF_NAME)) {
        console.log('MALFORMED plugin name: ' + pluginName)
        return ''
    }
    return pluginName.substring(START_OF_NAME.length)
}

function darwinExec(command: string, ...args: string[]) : ProgressMessage {
    const result = {message: '', details: '', error: false }
    try {
        const syncResult = spawnSync(command, args, {stdio: 'pipe', encoding: 'utf8'})
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
    return darwinExec('mkdir', '-p', dir)
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

// convenience wrapper to darwinExec for a copy command
function darwinCopy(src, dst: string): ProgressMessage {
    return darwinExec('cp', src, dst)
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
    const result = darwinExec('tar', 'xzvf', tarballInfo.fullPath, '-C', tarballInfo.dir)
    if (!result.error) {
        result.message =  'Successful untar of ' + tarballInfo.fullPath
    }
    result.stepComplete = true
    return result
}

function darwinEnvPathContains(envPath, target) {
    return utils.stringContains(':' + envPath + ':', ':' + target + ':')
}

module.exports.steps = darwinSteps
module.exports.preinstall = preinstallDarwin

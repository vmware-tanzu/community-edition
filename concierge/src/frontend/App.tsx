import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { ipcRenderer } from 'electron';

// App imports
import Concierge from './components/Concierge';
import { AvailableInstallation, PreInstallation } from '../models/installation';

// App css imports
import '@cds/core/global.css'; // pre-minified version breaks
import '@cds/city/css/bundles/default.min.css';
import '@cds/core/global.min.css';
import '@cds/core/styles/theme.dark.min.css';
import './index.scss';
import { HashRouter } from 'react-router-dom';
import { ProgressMessage } from '_models/progressMessage';

let setPreinstallationData      // fxn (to set pre-installation data) that is set by Concierge via callback
let progressReceiver            // fxn (to receive a progress message from the backend) that is set by Concierge via callback
// ******************************************
// callbacks made available to Concierge
//
function refreshInstallations() {
    console.log('Sending app:pre-install-tanzu');
    ipcRenderer.send('app:pre-install-tanzu');
}

function install(installation: AvailableInstallation) {
    console.log('Sending app:install-tanzu message');
    ipcRenderer.send('app:install-tanzu', installation);
}

// This callback sends us a fxn we can use to set the preinstallation data (when that data arrives from the main process)
function registerSetPreinstallationFunction(setPreinstallationFunction: (preinstallation: PreInstallation) => void) {
    console.log('APP: SetPreinstallationFunction has been registered')
    setPreinstallationData = setPreinstallationFunction
}

function registerProgressReceiver(pr: (pm: ProgressMessage) => void) {
    console.log('APP: progress receiver has been registered')
    progressReceiver = pr
}

function runTanzu() {
    console.log('Sending app:launch-kickstart message');
    ipcRenderer.send('app:launch-kickstart');
/*
    console.log('Sending app:launch-tanzu-ui message');
    ipcRenderer.send('app:launch-tanzu-ui');
*/
}
//
// callbacks made available to Concierge
// ******************************************

ipcRenderer.on('app:pre-install-tanzu', (event, message) => {
    console.log('ipcRenderer got message pre-install-tanzu: ' + JSON.stringify(message))
    if (message && setPreinstallationData) {
        setPreinstallationData(message)
        return
    }
    if (!message) {
        console.log('Received pre-install-tanzu message, but there was no data!')
    }
    if (!setPreinstallationData) {
        console.log('Received pre-install-tanzu message, but unable to set preinstallation data cuz no fxn available')
    }
});

ReactDOM.render(
    <React.StrictMode>
        <HashRouter>
            <Concierge
                refreshInstallations={refreshInstallations}
                install={install}
                registerSetPreinstallationFunction={registerSetPreinstallationFunction}
                registerProgressReceiver={registerProgressReceiver}
                runTanzu={runTanzu}
            />
        </HashRouter>
    </React.StrictMode>,
    document.getElementById('concierge')
);
/*

ipcRenderer.on('app:pre-install-tanzu', (event, message) => {
    console.log('ipcRenderer got message pre-install-tanzu: ' + JSON.stringify(message))
    if (message) {
        let canInstall = true
        const preInstall = message as PreInstallation
        const existingInstallation = preInstall.existingInstallation
        let displayText = ''
        if (existingInstallation && existingInstallation.path) {
            const displayEdition = existingInstallation.edition ? existingInstallation.edition : ''
            const binaryVersion = existingInstallation.tanzuBinaryVersion ? 'tanzu binary ' + existingInstallation.tanzuBinaryVersion :
                'unknown tanzu binary version'
            const editionVersion = existingInstallation.editionVersion
                ? existingInstallation.edition + ' ' + existingInstallation.editionVersion
                : 'unknown ' + displayEdition + ' version using ' + binaryVersion
            displayText = 'We notice you have an existing ' + displayEdition + ' installation at ' +
                existingInstallation.path + ' (' + editionVersion + ')'
            if (!canInstallOver(existingInstallation.edition)) {
                // TODO: handle non-TCE versions in this message
                displayText += '\n\nOh, no!\nWe are not sophisticated enough to install TCE over this existing ' + displayEdition
                    + ' installation. So sorry.\nYou can try removing the ' + displayEdition +
                    ' installation manually, and then try us again, or '
                displayText += 'you can create a new user account and then we should be able to install TCE if you log in as the new user.'
                displayText += '\n\nFor more information, see http://helpful.url.here'
                canInstall = false
            }
        } else {
            displayText = 'No existing installation detected; that makes things easier.'
        }

        if (canInstall) {
            if (!preInstall.availableInstallations || preInstall.availableInstallations.length === 0) {
                displayText += '\n\nHorrors! No available installations were detected. I won\'t be able to install anything.' +
                    ' Major fail!\n\n(Note: the Concierge distribution may be defective, or the tarball may have been manually deleted.' +
                    ` I was looking for installation tarballs in these directories:\n${preInstall.dirInstallationTarballsExpected.join('\n')}.)`
            } else if (preInstall.availableInstallations.length === 1) {
                const onlyInstall = preInstall.availableInstallations[0]
                chosenInstallation = onlyInstall
                displayText += '\n\nAre you ready to install Tanzu (' + chosenInstallation.edition.toUpperCase() + ' ' + chosenInstallation.version + ')?'
            } else {
                displayText += '\nApparently we have ' + preInstall.availableInstallations.length + ' installations to choose from!\n'
                preInstall.availableInstallations.every(availInstall => displayText += `(${availInstall.edition.toUpperCase()} ${availInstall.version})\n`)
                chosenInstallation = preInstall.availableInstallations[preInstall.availableInstallations.length-1]
                displayText += '\n\nAre you ready to install Tanzu (' + chosenInstallation.edition.toUpperCase() + ' ' + chosenInstallation.version + ')?'
            }
            console.log(`PREINSTALL: looking for installation tarballs in these directories:\n${preInstall.dirInstallationTarballsExpected.join('\n')}`)
        }
        const ctlInfo = document.getElementById('existingTanzuInfo')
        if (ctlInfo) {
            ctlInfo.innerText = displayText
        }
    }
});
*/

function canInstallOver(existingEdition: string): boolean {
    return !existingEdition || existingEdition === 'TCE' || existingEdition === 'BOZO'
}

ipcRenderer.on('app:install-progress', (event, progressMessageObject) => {
    if (progressReceiver) {
        progressReceiver(progressMessageObject)
    } else {
        console.log('PROGRAMMER ERROR (APP): received progress message, but no receiver to give it to!')
    }

    if (progressMessageObject) {
        if (progressMessageObject.error) {
            addMessage(`--- ERROR ---\nSTEP: ${progressMessageObject.step} MSG: ${progressMessageObject.message} DETAILS: ${progressMessageObject.details}`)
        } else if (progressMessageObject.message || progressMessageObject.details){
            const step = progressMessageObject.step ? `STEP: ${JSON.stringify(progressMessageObject.step)}` : ''
            const details = progressMessageObject.details ? `DETAILS: ${JSON.stringify(progressMessageObject.details)}` : ''
            const data = progressMessageObject.data ? `DATA: ${JSON.stringify(progressMessageObject.data)}` : ''
            console.log(`${step} MSG: ${progressMessageObject.message} ${details} ${data}`)
        }
        if (progressMessageObject.percentComplete) {
            displayPercentage(progressMessageObject.percentComplete + '%')
        }
        if (progressMessageObject.stepStarting) {
            displayStep('STARTING: ' + progressMessageObject.message)
            displayPercentage('')
        }
        if (progressMessageObject.stepComplete) {
            displayStep('COMPLETE: ' + progressMessageObject.message)
            displayPercentage('')
        }
        if (progressMessageObject.installComplete) {
            if (progressMessageObject.error) {
                displayStep('INSTALLATION FAILED: ' + progressMessageObject.message)
            } else {
                displayStep('INSTALLATION SUCCEEDED: ' + progressMessageObject.message)
                postInstall()
            }
            displayPercentage('')
        }
        if (progressMessageObject.installStarting) {
            addMessage(progressMessageObject.message)
        }
    } else {
        console.log('ERROR: received app:install-progress message with no object!')
    }
});

ipcRenderer.on('app:plugin-list-response', (event, pluginList) => {
    const msg = `I see you now have these plugins installed: \n[${pluginList.join(']\n[')}]\nIsn\'t that nice?`
    displayStep(msg)
})

function addMessage(message: string) {
    const time = displayTime()
    const ctlProgressDisplay = document.getElementById('installProgressDisplay')
    if (ctlProgressDisplay) {
        ctlProgressDisplay.innerText = ctlProgressDisplay.innerText + time + ' > ' + message + '\n'
    }
}

function displayInElement(element, message: string) {
    const ctl = document.getElementById(element)
    if (ctl) {
        ctl.innerText = message
    }
}

function displayPercentage(message: string) {
    displayInElement('percentComplete', message)
}

function displayStep(message: string) {
    displayInElement('stepName', message)
}

function displayTime() {
    const now = new Date()
    const minutes = now.getMinutes() < 10 ? '0' + now.getMinutes() : '' + now.getMinutes()
    return now.getHours() + ":" + minutes + ":" + now.getSeconds()
}

function postInstall() {
    ipcRenderer.send('app:plugin-list-request')
}


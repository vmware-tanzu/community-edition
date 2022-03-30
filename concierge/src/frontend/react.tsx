import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { ipcRenderer } from 'electron';

import Concierge from './components/Concierge';
import { PreInstallation } from '../models/installation';

let chosenInstallation
ReactDOM.render(<Concierge />, document.getElementById('concierge'));

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
            if (!preInstall.availableInstallations) {
                displayText += 'Horrors! No available installations were detected. Somehow I\'m missing the ability to install anything.' +
                    ' Major fail!'
            } else if (preInstall.availableInstallations.length === 1) {
                const onlyInstall = preInstall.availableInstallations[0]
                chosenInstallation = onlyInstall
                displayText += '\n\nAre you ready to install Tanzu (' + chosenInstallation.edition.toUpperCase() + ' ' + chosenInstallation.version + ')?'
            } else {
                displayText += '\nApparently we have ' + preInstall.availableInstallations.length + ' installations to choose from!' +
                    '\n' + JSON.stringify(preInstall.availableInstallations)
                chosenInstallation = preInstall.availableInstallations[preInstall.availableInstallations.length-1]
                displayText += '\n\nAre you ready to install Tanzu (' + chosenInstallation.edition.toUpperCase() + ' ' + chosenInstallation.version + ')?'
            }
        }
        document.getElementById('existingTanzuInfo').innerText = displayText
    }
});

function canInstallOver(existingEdition: string): boolean {
    return !existingEdition || existingEdition === 'TCE' || existingEdition === 'BOZO'
}

ipcRenderer.on('app:install-progress', (event, progressMessageObject) => {
    if (progressMessageObject) {
        const time = displayTime()
        const ctlProgressDisplay = document.getElementById('installProgressDisplay')
        const currentProgress = ctlProgressDisplay.innerText
        let messageToAdd = ''
        if (progressMessageObject.error) {
            messageToAdd = '--- ERROR ---\n'
        }
        messageToAdd += time + ' > '
        if (progressMessageObject.step) {
            messageToAdd += 'STEP: ' + progressMessageObject.step + ' '
        }
        if (progressMessageObject.message) {
            messageToAdd += 'MSG: ' + progressMessageObject.message + '\n'
        }
        if (progressMessageObject.details) {
            messageToAdd += 'DETAILS: ' + progressMessageObject.details + '\n'
        }
        if (progressMessageObject.percentComplete) {
            const ctlPercentComplete = document.getElementById('percentComplete')
            ctlPercentComplete.innerText = progressMessageObject.percentComplete + '%'
        }
        if (messageToAdd) {
            ctlProgressDisplay.innerText = currentProgress + messageToAdd + '\n'
        }
    }
});

function displayTime() {
    const now = new Date()
    const minutes = now.getMinutes() < 10 ? '0' + now.getMinutes() : '' + now.getMinutes()
    return now.getHours() + ":" + minutes + ":" + now.getSeconds()
}

const installButton = document.getElementById('buttonInstall');
if (installButton) {
    installButton.addEventListener('click', function () {
        console.log('Sending app:install-tanzu message');
        ipcRenderer.send('app:install-tanzu', chosenInstallation); // ipcRender.send will pass the information to main process
    });
}

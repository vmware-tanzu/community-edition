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
        document.getElementById('existingTanzuInfo').innerText = displayText
    }
});

function canInstallOver(existingEdition: string): boolean {
    return !existingEdition || existingEdition === 'TCE' || existingEdition === 'BOZO'
}

ipcRenderer.on('app:install-progress', (event, progressMessageObject) => {
    if (progressMessageObject) {
        let messageToAdd = ''
        if (progressMessageObject.error) {
            messageToAdd = '--- ERROR ---\n'
        }
        if (progressMessageObject.step) {
            messageToAdd += 'STEP: ' + progressMessageObject.step + ' '
        }
        if (progressMessageObject.message) {
            messageToAdd += 'MSG: ' + progressMessageObject.message + '\n'
        }
        if (progressMessageObject.details) {
            // messageToAdd += 'DETAILS: ' + progressMessageObject.details + '\n'
            console.log(`STEP: ${progressMessageObject.step} sez ${progressMessageObject.details}`)
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
            }
            displayPercentage('')
        }
        if (progressMessageObject.installStarting) {
            messageToAdd = progressMessageObject.message
        }

        if (messageToAdd) {
            addMessage(messageToAdd)
        }
    }
});

function addMessage(message: string) {
    const time = displayTime()
    const ctlProgressDisplay = document.getElementById('installProgressDisplay')
    ctlProgressDisplay.innerText = ctlProgressDisplay.innerText + time + ' > ' + message + '\n'
}

function displayInElement(element, message: string) {
    document.getElementById(element).innerText = message
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

const installButton = document.getElementById('buttonInstall');
if (installButton) {
    installButton.addEventListener('click', function () {
        console.log('Sending app:install-tanzu message');
        ipcRenderer.send('app:install-tanzu', chosenInstallation); // ipcRender.send will pass the information to main process
    });
}
const preinstallButton = document.getElementById('buttonPreInstall');
if (preinstallButton) {
    preinstallButton.addEventListener('click', function () {
        console.log('Sending app:pre-install-tanzu message');
        ipcRenderer.send('app:pre-install-tanzu');
    });
}

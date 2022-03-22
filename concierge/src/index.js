import React from 'react'
import { render } from 'react-dom'
import { ipcRenderer } from 'electron';

import Concierge from './components/Concierge';

// Since we are using HtmlWebpackPlugin WITHOUT a template, we should create our own root node in the body element before rendering into it
let root = document.createElement('div')

root.id = 'root'
document.body.appendChild(root)

// Now we can render our application into it
render(<Concierge />, document.getElementById('root'))

ipcRenderer.on('app:existing-install-tanzu', (event, message) => {
    console.log('ipcRenderer got message existing-install-tanzu: ' + JSON.stringify(message))
    if (message && message.path && message.version) {
        document.getElementById('existingTanzuInfo').innerText = 'We notice you have an existing ' + message.edition + ' installation: version ' + message.version +
            ' installed to ' + message.path
    }
});

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
        ipcRenderer.send('app:install-tanzu', 'test'); // ipcRender.send will pass the information to main process
    });
}

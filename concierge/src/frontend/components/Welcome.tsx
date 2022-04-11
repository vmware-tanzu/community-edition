import { CdsButton } from '@cds/react/button';
import { CdsSelect } from '@cds/react/select';
import * as React from 'react';
import { useNavigate } from 'react-router-dom';

import SelectInstallation from '_frontend/components/SelectInstallation';
import { ipcRenderer } from 'electron';

const Name = React.createContext<string[]>(['a','b']);

function onSelectInstallation(ev) {
    console.log('SHIMON SEZ: Got me a select installation event')
    console.log(`SHIMON SEZ: event target attributes: ${JSON.stringify(ev.target.attributes)}`)
    console.log(`SHIMON SEZ: event target id: ${ev.target.id}`)
    console.log(`SHIMON SEZ: event target value: ${ev.target.value}`)
}

function Welcome(props) {
    let navigate = useNavigate();
    let [availInstalls, setAvailInstalls] = React.useState(['1', '2', '3'])
    ipcRenderer.on('app:pre-install-tanzu', (event, message) => {
        console.log(`Got message ${JSON.stringify(message)}`)
        setAvailInstalls(message.availableInstallations.map(installation => installation.version))
    })
    return (
        <div>
        <h1>Welcome to the Tanzu Concierge!</h1>

    <p>We're here to help you install the Tanzu Commandline Interface (CLI) and the relevant plugins
    to create and manage clusters, install community packages and more.</p>
    <div id="existingTanzuInfo"></div>

    <SelectInstallation availableInstallations={availInstalls} onSelect={onSelectInstallation} refreshInstallations={props.refreshInstallations} />

    <p>&nbsp;</p>
    <button type="button" id="buttonInstall">
        Install Tanzu, baby!
    </button>
    <div id="stepName"></div><div id="percentComplete"></div>

    <p>&nbsp;</p>
    <div id="installProgressDisplay"></div>

    <p>&nbsp;</p>

    <button type="button" id="buttonLaunchKickstart">
        Launch Kickstart, baby!
    </button>

    <p>&nbsp;</p>
    <button type="button" id="buttonLaunchTanzuUi">
        Launch Tanzu UI, baby!
    </button>
    <CdsButton onClick={() => {
        navigate("/install", {})
    }}>Go, Go, Go</CdsButton>

        </div>

);
}

export default Welcome

import { CdsButton } from '@cds/react/button';
import * as React from 'react';
import { useNavigate } from 'react-router-dom';

import SelectInstallation from '_frontend/components/SelectInstallation';
import { AvailableInstallation } from '_models/installation';
import ExistingInstallation from '_frontend/components/ExistingInstallation';

// Expects:
// props.preinstallationData            a PreInstallation object describing available installation versions and existing installation (if any)
// props.selectedInstallationVersion    the installation version the user has selected (if any)
// props.selectedInstallation           the installation object associated with the user's selected version (if any)
// props.refreshInstallations           a function to call which will refresh the installation list in the PreInstallation object
// props.setInstallationVersion         a function to call which sets the installation version (when the user selects a version)
// props.runTanzu                       a function to call which will run the Tanzu UI (given an existing installation)
// props.install                        a function to call which will install the Tanzu binaries associated with the passed parameter
function Welcome(props) {
    validateProps(props)
    let navigate = useNavigate();
    let availInstalls = [] as AvailableInstallation[]
    if (props.preinstallationData) {
        availInstalls = props.preinstallationData.availableInstallations
    } else {
        console.log('WELCOME has NO preinstallation data')
    }

    const availableInstallationVersions = availInstalls.map<string>(installation => installation.version)
    return (
        <div>
        <h1>Welcome to the Tanzu Concierge!</h1>

    <p>We're here to help you install the Tanzu Commandline Interface (CLI) and the relevant plugins
    to create and manage clusters, install community packages and more.</p>

    <ExistingInstallation
        existingInstallation={props.preinstallationData?.existingInstallation}
        runTanzu={props.runTanzu}
    />

    <SelectInstallation
        availableInstallationVersions={availableInstallationVersions}
        onSelect={props.setInstallationVersion}
        refreshInstallations={props.refreshInstallations}
    />
    <CdsButton onClick={() => {
        props.install(props.selectedInstallation)
        navigate("/install", {})
    }}>INSTALL TANZU ({props.selectedInstallationVersion ? 'enabled' : 'disabled'})</CdsButton>

    <p>&nbsp;</p>
    <div id="existingTanzuInfo"></div>
    <div id="stepName"></div><div id="percentComplete"></div>

    <p>&nbsp;</p>
    <div id="installProgressDisplay"></div>


        </div>

);
}

function validateProps(props) {
    // TODO: more validations
    if (!props.install) {
        console.log('PROGRAMMER ERROR: (ExistingInstallation) props.install is undefined!')
    }
}
export default Welcome

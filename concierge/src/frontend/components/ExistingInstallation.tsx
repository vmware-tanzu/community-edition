import * as React from 'react';
import SelectBox from '_frontend/components/SelectBox';
import { CdsButton } from '@cds/react/button';

// Expects:
// props.existingInstallation           InstallationVersion object describing installed version (if there is one)
// props.runTanzu                       callback will run the current version

function ExistingInstallation(props) {
    validateProps(props)
    const hasVersion = props.existingInstallation && props.existingInstallation.editionVersion && props.existingInstallation.editionVersion.length > 0
    return (
        <div>
            <p>Installed version: {hasVersion ? props.existingInstallation.editionVersion : 'None'}</p>
            <CdsButton onClick={() => {
                props.runTanzu()
            }}>OPEN TANZU UI ({hasVersion ? 'enabled' : 'disabled'})</CdsButton>
        </div>);
}

function validateProps(props) {
    if (!props.runTanzu) {
        console.log('PROGRAMMER ERROR: (ExistingInstallation) props.runTanzu is undefined!')
    }
}


export default ExistingInstallation

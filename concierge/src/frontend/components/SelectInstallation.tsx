import * as React from 'react';
import SelectBox from '_frontend/components/SelectBox';
import { CdsButton } from '@cds/react/button';

// Expects:
// props.availableInstallations     list of strings identifying versions
// props.onSelect                   callback that takes an event (use event.target.value to get selected install version)
// props.refreshInstallations       callback that will refresh the installation list
function SelectInstallation(props) {
    return (
        <div>
            <h1>Select from these installations <CdsButton onClick={() => {props.refreshInstallations()}}>Refresh</CdsButton></h1>

            <SelectBox data={props.availableInstallations} onSelect={props.onSelect} label="Installation" id="installationId" />
        </div>);
}


export default SelectInstallation

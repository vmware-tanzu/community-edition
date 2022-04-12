import * as React from 'react';
import { useNavigate } from 'react-router-dom';
import { ProgressMessage } from '_models/progressMessage';
import { CdsButton } from '@cds/react/button';

function displayProgressMessage(pm: ProgressMessage): string {
    const stepDisplay = pm.step ? `${pm.step}: ` : ''
    if (pm.message && pm.details) {
        return `${stepDisplay}${pm.message} (${pm.details})`
    }
    if (pm.details) {
        return `${stepDisplay}(${pm.details})`
    }
    if (pm.message) {
        return `${stepDisplay}${pm.message}`
    }
    if (pm.step) {
        return `STEP: ${pm.step}`
    }
    return `Unrecognized message: ${JSON.stringify(pm)}`
}

function Install(props) {
    const navigate = useNavigate();
    return (
        <div>
            <h1>Installing Tanzu</h1>
            <div>
                <h3>Installation Messages</h3>
                <ul>
                {props.progressMessages.map(pm => <li>{displayProgressMessage(pm)}</li>)}
                </ul>
            </div>
            <CdsButton onClick={() => {
                navigate("/", {})
            }}>HOME (back)</CdsButton>
        </div>);
}

export default Install

import * as React from 'react';
import { ProgressMessage } from '_models/progressMessage';

function displayProgressMessage(pm: ProgressMessage): string {
    const stepDisplay = pm.step ? `${pm.step}: ` : ''
    if (pm.message && pm.details) {
        return `${stepDisplay}${pm.message} (${pm.details})`
    }
    if (pm.details) {
        return `${stepDisplay}(${pm.details})`
    }
    return `STEP: ${pm.step}`
}

function Install(props) {
    return (
        <div>
            <h1>Install Page!</h1>
            <div>
                <h3>Installation Messages</h3>
                <ul>
                {props.progressMessages.map(pm => <li>{displayProgressMessage(pm)}</li>)}
                </ul>
            </div>
        </div>);
}

export default Install

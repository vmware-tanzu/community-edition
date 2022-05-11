// React imports
import { ChangeEvent } from 'react';
// App imports
import { ManagementCluster } from '../../swagger-api';

// Utility method for retrieving a value from a change event.
// Returns true/false from checkboxes; otherwise the value in the event object
export function getValueFromChangeEvent(evt: ChangeEvent<HTMLSelectElement>) {
    let value
    if (evt.target.type === 'checkbox') {
        // @ts-ignore
        value = evt.target['checked']
    } else {
        value = evt.target.value
    }
    return value
}

export function getSelectedManagementCluster(state: any): ManagementCluster {
    return state.data.SELECTED_MANAGEMENT_CLUSTER
}


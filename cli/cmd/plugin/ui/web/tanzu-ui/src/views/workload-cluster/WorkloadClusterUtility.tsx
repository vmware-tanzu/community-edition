// React imports
import { ChangeEvent } from 'react';
// App imports
import { ManagementCluster } from '../../swagger-api';
import { STORE_SECTION_FORM } from '../../state-management/reducers/Form.reducer';

// Utility method for retrieving a value from a change event.
// Returns true/false from checkboxes; otherwise the value in the event object
export function getValueFromChangeEvent(evt: ChangeEvent<HTMLSelectElement>) {
    let value;
    if (evt.target.type === 'checkbox') {
        value = evt.target['checked'];
    } else {
        value = evt.target.value;
    }
    return value;
}

export function getSelectedManagementCluster(state: any): ManagementCluster {
    return state[STORE_SECTION_FORM].SELECTED_MANAGEMENT_CLUSTER;
}

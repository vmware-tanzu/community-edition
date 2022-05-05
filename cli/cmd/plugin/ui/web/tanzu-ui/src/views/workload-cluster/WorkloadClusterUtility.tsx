import { ManagementCluster } from '../../shared/models/ManagementCluster';
import { ChangeEvent } from 'react';

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

// NOTE: Rather than store the CC variables at the top level of our data store,
//       we accumulate them in a single object per management cluster.
//       This makes using them much easier (later), and prevents "collisions" if the user switches management clusters.
function getClusterVariableData(managementCluster: ManagementCluster, state: any): any {
    if (!managementCluster || !managementCluster.name) {
        console.error('getClusterVariableData() called with undefined/unnamed cluster')
        return {}
    }
    const ccVarData = { ...state.data.CLUSTER_CLASS_VARIABLE_VALUES }
    return ccVarData[managementCluster.name] || {}
}

export function keyClusterClassVariableData(): string {
    return 'CLUSTER_CLASS_VARIABLE_VALUES'
}

export function modifyClusterVariableDataItem(key: string, value: any, managementCluster: ManagementCluster, state: any): any {
    const ccVarClusterData = { ...getClusterVariableData(managementCluster, state) }
    if (key) {
        if (value) {
            ccVarClusterData[key] = value
        } else {
            // NOTE: A boolean field with a value of FALSE gets omitted. That works for all our current fields.
            delete ccVarClusterData[key]
        }
    } else {
        console.error(`setClusterVariableDataItem was called with a null key (MC=${managementCluster?.name})`)
    }
    return ccVarClusterData
}

// NOTE: we assume that state.data.CLUSTER_CLASS_VARIABLE_VALUES is never null or undefined
export function getSelectedManagementCluster(state: any): ManagementCluster {
    return state.data.SELECTED_MANAGEMENT_CLUSTER
}

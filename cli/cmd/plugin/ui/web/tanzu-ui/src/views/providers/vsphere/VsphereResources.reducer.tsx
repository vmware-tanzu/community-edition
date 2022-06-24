import { VsphereResourceAction } from '../../../shared/types/types';
import { ReducerDescriptor } from '../../../shared/utilities/Reducer.utils';
import { VSPHERE_ADD_RESOURCES, VSPHERE_DELETE_RESOURCES } from '../../../state-management/actions/Resources.actions';

export const STORE_SECTION_VSPHERE_RESOURCES = 'resources';

// The resources reducer builds up state by ensuring an object associated with action.datacenter,
// and then assigning object[action.resourceName] = action.payload.
// So if the action object were:
// { datacenter: dc-10, resourceName: osImages, payload: [obj1, obj2, obj3] }, then we would expect
// to see state['dc-10']['osImages'] === [obj1, obj2, obj3]
function vsphereResourcesReducer(state: any, action: VsphereResourceAction) {
    const newState = { ...state };
    if (!action.datacenter) {
        console.error(`vsphereResourcesReducer received action ${JSON.stringify(action)} which has no datacenter!`);
        return newState;
    }
    if (!action.resourceName) {
        console.error(`vsphereResourcesReducer received action ${JSON.stringify(action)} which has no resourceName!`);
        return newState;
    }
    if (action.type === VSPHERE_ADD_RESOURCES) {
        if (!newState[action.datacenter]) {
            newState[action.datacenter] = {};
        }
        newState[action.datacenter][action.resourceName] = action.payload;
    } else if (action.type === VSPHERE_DELETE_RESOURCES) {
        if (newState[action.datacenter]) {
            delete newState[action.datacenter][action.resourceName];
            if (Object.keys(newState[action.datacenter]).length === 0) {
                delete newState[action.datacenter];
            }
        }
    }
    console.log(`New resources state: ${JSON.stringify(newState)} after action ${JSON.stringify(action)}`);
    return newState;
}

export const vsphereResourceReducerDescriptor = {
    name: 'vsphere resource reducer',
    reducer: vsphereResourcesReducer,
    storeSection: STORE_SECTION_VSPHERE_RESOURCES,
    actionTypes: [VSPHERE_ADD_RESOURCES, VSPHERE_DELETE_RESOURCES],
} as ReducerDescriptor;

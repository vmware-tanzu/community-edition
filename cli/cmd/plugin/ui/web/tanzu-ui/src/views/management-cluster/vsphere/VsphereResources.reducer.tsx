import { Action, VsphereResourceAction } from '../../../shared/types/types';
import { ReducerDescriptor } from '../../../shared/utilities/Reducer.utils';
import { VSPHERE_ADD_RESOURCES, VSPHERE_DELETE_RESOURCES } from '../../../state-management/actions/Resources.actions';

// The resources reducer builds up state by ensuring an object associated with action.datacenter,
// and then assigning object[action.resourceName] = action.payload.
// So if the action object were:
// { datacenter: dc-10, resourceName: osImages, payload: [obj1, obj2, obj3] }, then we would expect
// to see state['dc-10']['osImages'] === [obj1, obj2, obj3]
function vsphereResourcesReducer(state: any, action: VsphereResourceAction) {
    const newState = { ...state };
    if (action.type === VSPHERE_ADD_RESOURCES) {
        if (!newState[action.datacenter]) {
            newState[action.datacenter] = {};
        }
        newState[action.datacenter][action.resourceName] = action.payload;
        console.log(`New resources state: ${JSON.stringify(newState)}`);
    }
    return newState;
}

export const vsphereResourceReducerDescriptor = {
    name: 'vsphere resource reducer',
    reducer: vsphereResourcesReducer,
    storeSection: 'resources',
    actionTypes: [VSPHERE_ADD_RESOURCES, VSPHERE_DELETE_RESOURCES],
} as ReducerDescriptor;

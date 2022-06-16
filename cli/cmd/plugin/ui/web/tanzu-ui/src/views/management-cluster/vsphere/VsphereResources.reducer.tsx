import { ADD_RESOURCES, DELETE_RESOURCES } from '../../../state-management/actions/Resources.actions';
import { Action, VsphereResourceAction } from '../../../shared/types/types';
import { ReducerDescriptor } from '../../../shared/utilities/Reducer.utils';

// The resources reducer builds up state by ensuring an object associated with action.datacenter,
// and then assigning object[action.resourceName] = action.payload.
// So if the action object were:
// { datacenter: dc-10, resourceName: osImages, payload: [obj1, obj2, obj3] }, then we would expect
// to see state['dc-10']['osImages'] === [obj1, obj2, obj3]
function vsphereResourcesReducer(state: any, action: Action) {
    const newState = { ...state };
    if (action.type === ADD_RESOURCES) {
        const resourceAction = action as VsphereResourceAction;
        if (!newState[resourceAction.datacenter]) {
            newState[resourceAction.datacenter] = {};
        }
        newState[resourceAction.datacenter][resourceAction.resourceName] = resourceAction.payload;
        console.log(`New resources state: ${JSON.stringify(newState)}`);
    }
    return newState;
}

export const vsphereResourceReducerDescriptor = {
    name: 'vsphere resource reducer',
    reducer: vsphereResourcesReducer,
    storeSection: 'resources',
    actionTypes: [ADD_RESOURCES, DELETE_RESOURCES],
} as ReducerDescriptor;

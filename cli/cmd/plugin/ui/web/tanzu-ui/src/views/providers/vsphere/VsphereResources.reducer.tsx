import { ResourceAction } from '../../../shared/types/types';
import { ReducerDescriptor } from '../../../shared/utilities/Reducer.utils';
import { RESOURCE } from '../../../state-management/actions/Resources.actions';

export const STORE_SECTION_VSPHERE_RESOURCES = 'resources';

// The resources reducer builds up state by ensuring an object associated with action.datacenter,
// and then assigning object[action.resourceName] = action.payload.
// So if the action object were:
// { datacenter: dc-10, resourceName: osImages, payload: [obj1, obj2, obj3] }, then we would expect
// to see state['dc-10']['osImages'] === [obj1, obj2, obj3]
function vsphereResourcesReducer(state: any, action: ResourceAction) {
    const newState = { ...state };
    if (!action.resourceName) {
        console.error(`vsphereResourcesReducer received action ${JSON.stringify(action)} which has no resourceName!`);
        return newState;
    }
    if (action.type === RESOURCE.VSPHERE_ADD_RESOURCES) {
        newState[action.resourceName] = action.payload;
    } else if (action.type === RESOURCE.VSPHERE_DELETE_RESOURCES) {
        delete newState[action.resourceName];
    }
    console.log(`New resources state: ${JSON.stringify(newState)} after action ${JSON.stringify(action)}`);
    return newState;
}

export function getResource(resourceName: string, store: any) {
    if (!store || !store[STORE_SECTION_VSPHERE_RESOURCES] || !resourceName) {
        return undefined;
    }
    return store[STORE_SECTION_VSPHERE_RESOURCES][resourceName];
}

export const vsphereResourceReducerDescriptor = {
    name: 'vsphere resource reducer',
    reducer: vsphereResourcesReducer,
    storeSection: STORE_SECTION_VSPHERE_RESOURCES,
    actionTypes: [RESOURCE.VSPHERE_ADD_RESOURCES, RESOURCE.VSPHERE_DELETE_RESOURCES],
} as ReducerDescriptor;

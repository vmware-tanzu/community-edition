import { AwsResourceAction } from '../../../shared/types/types';
import { ReducerDescriptor } from '../../../shared/utilities/Reducer.utils';
import { AWS_ADD_RESOURCES, AWS_DELETE_RESOURCES } from '../../../state-management/actions/Resources.actions';

export const STORE_SECTION_AWS_RESOURCES = 'resources';

// The resources reducer builds up state by ensuring an object associated with action.datacenter,
// and then assigning object[action.resourceName] = action.payload.
// So if the action object were:
// { datacenter: dc-10, resourceName: osImages, payload: [obj1, obj2, obj3] }, then we would expect
// to see state['dc-10']['osImages'] === [obj1, obj2, obj3]
function awsResourcesReducer(state: any, action: AwsResourceAction) {
    const newState = { ...state };
    if (!action.resourceName) {
        console.error(`awsResourcesReducer received action ${JSON.stringify(action)} which has no resourceName!`);
        return newState;
    }
    if (action.type === AWS_ADD_RESOURCES) {
        newState[action.resourceName] = action.payload;
    } else if (action.type === AWS_DELETE_RESOURCES) {
        delete newState[action.resourceName];
    }
    console.log(`New resources state: ${JSON.stringify(newState)} after action ${JSON.stringify(action)}`);
    return newState;
}

export function getResource(resourceName: string, store: any) {
    if (!store || !store[STORE_SECTION_AWS_RESOURCES] || !resourceName) {
        return undefined;
    }
    return store[STORE_SECTION_AWS_RESOURCES][resourceName];
}

export const awsResourceReducerDescriptor = {
    name: 'aws resource reducer',
    reducer: awsResourcesReducer,
    storeSection: STORE_SECTION_AWS_RESOURCES,
    actionTypes: [AWS_ADD_RESOURCES, AWS_DELETE_RESOURCES],
} as ReducerDescriptor;

import { Action } from '../../shared/types/types';
import { ADD_RESOURCES } from '../actions/Resources.actions';

// The resources reducer builds up state by ensuring an object associated with action.locationData,
// and then assigning object[action.field] = action.payload.
// So if the action object were:
// { locationData: dc-10, field: osImages, payload: [obj1, obj2, obj3] }, then we would expect
// to see state['dc-10']['osImages'] === [obj1, obj2, obj3]
export function resourcesReducer(state: any, action: Action) {
    const newState = { ...state };
    switch (action.type) {
        case ADD_RESOURCES:
            if (!newState[action.locationData]) {
                newState[action.locationData] = {};
            }
            newState[action.locationData][action.field] = action.payload;
            console.log(`New resources state: ${JSON.stringify(newState)}`);
            break;
    }
    return newState;
}

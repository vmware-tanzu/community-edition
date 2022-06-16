import { Action } from '../types/types';

interface ReducerGroup {
    name: string;
    reducers: ReducerDescriptor[];
}
export interface ReducerDescriptor {
    name: string;
    storeSection: string;
    reducer: (state: any, action: Action) => any;
    actionTypes: string[];
}
export const groupedReducers = (reducerGroup: ReducerGroup) => {
    const reducerGroupMapping = createReducerGroupMapping(reducerGroup);
    return (state: any, action: Action) => {
        const reducerDescriptor = reducerGroupMapping.get(action.type);
        if (!reducerDescriptor) {
            console.error(
                `Mapped reducer ${JSON.stringify(reducerGroup)} is unable to find a reducer for action ${JSON.stringify(action)}`
            );
            return state;
        }
        const newState = { ...state };
        const section = reducerDescriptor.storeSection;
        newState[section] = reducerDescriptor.reducer(newState[section], action);
        return newState;
    };
};

function createReducerGroupMapping(reducerGroup: ReducerGroup): Map<string, ReducerDescriptor> {
    // Our goal is to create a single map that associates an action type with a single reducer
    // To do that, we cycle through all the ReducerDescriptors in the group and for each one we add all their action types to our map.
    // If we detect a duplicate (that is, an action type claimed by two ReducerDescriptors, or within a single ReducerDescriptor), we
    // complain to the console.
    return reducerGroup.reducers.reduce<Map<string, ReducerDescriptor>>((mappedReducers, reducerDescriptor) => {
        reducerDescriptor.actionTypes.forEach((action) => {
            const existingReducerDescriptor = mappedReducers.get(action);
            if (existingReducerDescriptor) {
                reportDuplicateReducer(action, reducerGroup, reducerDescriptor, existingReducerDescriptor);
            } else {
                mappedReducers.set(action, reducerDescriptor);
            }
        });
        return mappedReducers;
    }, new Map<string, ReducerDescriptor>());
}

function reportDuplicateReducer(
    action: string,
    reducerGroup: ReducerGroup,
    newReducerDescriptor: ReducerDescriptor,
    oldReducerDescriptor: ReducerDescriptor
) {
    console.error(
        `Error while grouping reducers for action ${action}: reducer ${JSON.stringify(newReducerDescriptor)} ` +
            `finds that action already been claimed by ${JSON.stringify(oldReducerDescriptor)}. Group: ${JSON.stringify(reducerGroup)}`
    );
}

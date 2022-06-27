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
        const reducerDescriptorArray = reducerGroupMapping.get(action.type);
        if (!reducerDescriptorArray || reducerDescriptorArray.length === 0) {
            console.error(`Group reducer ${JSON.stringify(reducerGroup)} unable to find any reducers for action ${JSON.stringify(action)}`);
            return state;
        }
        return applyAllReducers(action, reducerDescriptorArray, state);
    };
};

function applyAllReducers(action: Action, reducerDescriptorArray: ReducerDescriptor[], state: any) {
    const newState = { ...state };
    reducerDescriptorArray.forEach((reducerDescriptor) => {
        const section = reducerDescriptor.storeSection;
        newState[section] = reducerDescriptor.reducer(newState[section], action);
    });
    return newState;
}

function createReducerGroupMapping(reducerGroup: ReducerGroup): Map<string, ReducerDescriptor[]> {
    reportDuplicateReducersIfNec(reducerGroup);
    // Our goal is to create a single map that associates an action type with an array of reducers that handle that action.
    // To do that, we cycle through all the ReducerDescriptors in the ReducerGroup. For each ReducerDescriptor,
    // we ensure each of its action types is in our map, and that the given reducer is in the array of reducers
    // that want to handle that action type.
    // When we detect that there are two ReducerDescriptors that want to handle a single action type, we currently
    // report the occurrence to the console (because currently it would be an error); in the future, we may want to remove
    // that reporting to the console (because it's legal for two reducers to respond to the same action).
    // However, if a single reducer lists the SAME action multiple times in its action types, we ignore the duplicates and report
    // that to the console; that is an error (albeit a benign one).
    return reducerGroup.reducers.reduce<Map<string, ReducerDescriptor[]>>((mappedReducers, reducerDescriptor) => {
        reducerDescriptor.actionTypes.forEach((action) => {
            const existingReducerDescriptorArray = mappedReducers.get(action) || [];
            if (existingReducerDescriptorArray.length > 0) {
                if (existingReducerDescriptorArray.includes(reducerDescriptor)) {
                    reportDuplicateAction(action, reducerDescriptor);
                    return;
                }
                reportMultipleReducers(action, reducerGroup, reducerDescriptor, existingReducerDescriptorArray);
            }
            const newReducerDescriptorArray = [...existingReducerDescriptorArray, reducerDescriptor];
            mappedReducers.set(action, newReducerDescriptorArray);
        });
        return mappedReducers;
    }, new Map<string, ReducerDescriptor[]>());
}

function reportMultipleReducers(
    action: string,
    reducerGroup: ReducerGroup,
    newReducerDescriptor: ReducerDescriptor,
    existingReducerDescriptors: ReducerDescriptor[]
) {
    const nExistingDescriptors = existingReducerDescriptors.length;
    const describeExisting = nExistingDescriptors === 1 ? 'one other reducer' : `${nExistingDescriptors} other reducers`;
    console.warn(
        `While grouping reducers for action ${action}: reducer ${JSON.stringify(newReducerDescriptor)} ` +
            `finds that action already been claimed by ${describeExisting}. Group: ${JSON.stringify(reducerGroup)}`
    );
}

function reportDuplicateAction(action: string, reducerDescriptor: ReducerDescriptor) {
    console.warn(
        `While grouping reducers for action ${action}: reducer ${JSON.stringify(reducerDescriptor)} ` +
            'either lists that action more than once, or the reducer itself is listed more than once in the group (duplicate ignored)'
    );
}

function reportDuplicateReducersIfNec(reducerGroup: ReducerGroup) {
    const nReducersInGroup = reducerGroup.reducers.length;
    const nUniqueReducersInGroup = new Set(reducerGroup.reducers).size;

    if (nReducersInGroup > nUniqueReducersInGroup) {
        console.warn(`Duplicate reducers found in reducer group ${JSON.stringify(reducerGroup)}`);
    }
}

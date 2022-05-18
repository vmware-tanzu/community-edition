import { Action } from '../types/types';

interface Reducers {
    [key: string]: (state: any, action: Action) => any;
}

export const combineReducers = (reducers: Reducers) => {
    return (state: any, action: Action) => {
        const newState: any = { ...state };
        for (const key in reducers) {
            newState[key] = reducers[key](state[key], action);
        }
        return newState;
    };
};

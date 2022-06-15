// App imports
import { APP_ENV_CHANGE, APP_ROUTE_CHANGE } from '../actions/App.actions';
import { Action } from '../../shared/types/types';
import { ReducerDescriptor } from '../../shared/utilities/Reducer.utils';

interface AppState {
    appEnv?: string;
    appRoute?: string;
}

function appReducer(state: AppState, action: Action) {
    let newState = { ...state };
    switch (action.type) {
        case APP_ENV_CHANGE:
            newState = {
                ...state,
                appEnv: action.payload.value,
            };
            break;
        case APP_ROUTE_CHANGE:
            newState = {
                ...state,
                appRoute: action.payload.value,
            };
            break;
    }
    return newState;
}

export const appReducerDescriptor = {
    name: 'app reducer',
    reducer: appReducer,
    storeSection: 'app',
    actionTypes: [APP_ENV_CHANGE, APP_ROUTE_CHANGE],
} as ReducerDescriptor;

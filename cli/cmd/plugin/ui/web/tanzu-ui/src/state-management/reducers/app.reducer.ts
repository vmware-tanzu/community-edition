// App imports
import { APP_ENV_CHANGE } from '../actions/app.actions';
import { Action } from '../../types/types';

interface AppState {
    appEnv?: string
}

export function appReducer (state: AppState, action: Action) {
    let newState = { ...state };
    switch (action.type) {
        case APP_ENV_CHANGE:
            newState =  {
                [action.payload.name]: action.payload.value
            };
    }
    return newState;
}
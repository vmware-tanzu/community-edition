// App imports
import { APP_ENV_CHANGE } from '../actions/App.actions';
import { Action } from '../../shared/types/types';

interface AppState {
    appEnv?: string;
}

export function appReducer(state: AppState, action: Action) {
    let newState = { ...state };
    switch (action.type) {
        case APP_ENV_CHANGE:
            newState = {
                [action.payload.name]: action.payload.value,
            };
    }
    return newState;
}

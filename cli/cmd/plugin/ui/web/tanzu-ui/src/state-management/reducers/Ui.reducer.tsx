// App imports
import { TOGGLE_NAV } from '../actions/Ui.actions';
import { Action } from '../../shared/types/types';

interface UIState {
    navExpanded: boolean
}

export function uiReducer (state: UIState, action: Action) {
    let newState = { ...state };
    switch (action.type) {
    case TOGGLE_NAV:
        newState =  {
            navExpanded: !state.navExpanded
        };
    }
    return newState;
}
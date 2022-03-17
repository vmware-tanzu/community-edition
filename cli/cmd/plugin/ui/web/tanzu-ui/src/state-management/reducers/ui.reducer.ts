// App imports
import { TOGGLE_NAV } from '../actions/actionTypes';
import { Action } from '../../types/types';

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
import { TOGGLE_NAV } from '../constants/actionTypes';
import { Action } from '../types/types';

interface UiState {
    navExpanded: boolean
}

export function uiReducer (state: UiState, action: Action) {
    let newState = { ...state };
    switch (action.type) {
    case TOGGLE_NAV:
        newState =  {
            navExpanded: !state.navExpanded
        };
    }
    return newState;
}
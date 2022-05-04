// App imports
import { TOGGLE_NAV, TOGGLE_WC_CC_OPTIONAL, TOGGLE_WC_CC_REQUIRED } from '../actions/Ui.actions';
import { Action } from '../../shared/types/types';

interface UIState {
    navExpanded: boolean,
    wcCcRequiredExpanded: boolean,
    wcCcOptionalExpanded: boolean,
}

export function uiReducer (state: UIState, action: Action) {
    let newState = { ...state };
    switch (action.type) {
    case TOGGLE_NAV:
        newState['navExpanded'] = !state.navExpanded
        break
    case TOGGLE_WC_CC_REQUIRED:
        newState['wcCcRequiredExpanded'] = !state.wcCcRequiredExpanded
        break
    case TOGGLE_WC_CC_OPTIONAL:
        newState['wcCcOptionalExpanded'] = !state.wcCcOptionalExpanded
        break
    }
    return newState;
}

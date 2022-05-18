// App imports
import { TOGGLE_APP_STATUS, TOGGLE_NAV, TOGGLE_WC_CC_CATEGORY } from '../actions/Ui.actions';
import { Action } from '../../shared/types/types';

interface UIState {
    isDeployInProgress: boolean;
    navExpanded: boolean;
    wcCcCategoryExpanded: { [category: string]: boolean };
}

export function uiReducer(state: UIState, action: Action) {
    const newState = { ...state };
    switch (action.type) {
        case TOGGLE_APP_STATUS:
            newState['isDeployInProgress'] = !state.isDeployInProgress;
            break;
        case TOGGLE_NAV:
            newState['navExpanded'] = !state.navExpanded;
            break;
        case TOGGLE_WC_CC_CATEGORY:
            newState['wcCcCategoryExpanded'] = createStateToggleCategory(state.wcCcCategoryExpanded, action.locationData);
            break;
    }
    console.log('APP UI:', newState);
    return newState;
}

// given an old categoryExpanded object, create a new categoryExpanded object (with the category toggled)
function createStateToggleCategory(oldCategoryObject: any, category: string): any {
    const oldToggleValue = oldCategoryObject[category] || false;
    return { ...oldCategoryObject, [category]: !oldToggleValue };
}

import { SelectionType } from './TreeSelect.interface';

/**
 * State
 */
export interface TreeSelectState {
    checked: string[];
    expanded: string[];
    selectionType: SelectionType;
}

export const initialState: TreeSelectState = {
    checked: [],
    expanded: [],
    selectionType: SelectionType.Single,
};

/**
 * Actions
 */
export const enum TreeSelectActions {
    UpdateChecked = 'updateChecked',
    UpdateExpanded = 'updateExpanded',
}

/**
 * Reducer for the TreeSelect state.
 * @param state
 * @param action
 * @returns
 */

export const treeSelectReducer = (state: TreeSelectState, action: { payload?: any; type: TreeSelectActions }): TreeSelectState => {
    switch (action.type) {
        case TreeSelectActions.UpdateChecked: {
            let newChecked = [];
            if (state.selectionType === SelectionType.Single) {
                const alreadyChecked = state.checked.includes(action.payload);
                newChecked = alreadyChecked ? [] : [action.payload];
            } else {
                const remaining = state.checked.filter((item) => item !== action.payload);
                const alreadyChecked = state.checked.includes(action.payload);
                newChecked = alreadyChecked ? [...remaining] : [...remaining, action.payload];
            }

            return {
                ...state,
                checked: [...newChecked],
            };
        }

        case TreeSelectActions.UpdateExpanded: {
            const remaining = state.expanded.filter((item) => item !== action.payload);
            const alreadyExpanded = state.expanded.find((item) => item === action.payload);
            const newExpanded = alreadyExpanded ? [...remaining] : [...remaining, action.payload];

            return {
                ...state,
                expanded: [...newExpanded],
            };
        }
    }
};

/**
 * State
 */
export interface DataAccordionState {
    active: boolean;
}

export const initialState: DataAccordionState = {
    active: false,
};

/**
 * Actions
 */
export const enum DataAccordionActions {
    ToggleAccordion,
}

/**
 * Reducer
 * @param state
 * @param action
 * @returns
 */
export const accordionReducer = (
    state: DataAccordionState,
    action: { payload?: number; type: DataAccordionActions }
): DataAccordionState => {
    switch (action.type) {
        case DataAccordionActions.ToggleAccordion:
            return {
                ...state,
                active: !state.active,
            };
    }
};

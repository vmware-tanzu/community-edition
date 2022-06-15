export interface DataAccordionState {
    active: boolean;
}
export const initialState: DataAccordionState = {
    active: false,
};

export const enum DataAccordionActions {
    ToggleAccordion,
}

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

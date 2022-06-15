// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import { STORE_SECTION_FORM } from '../reducers/Form.reducer';
import { STORE_SECTION_UI } from '../reducers/Ui.reducer';
import { StoreDispatch } from '../../shared/types/types';
import wcReducer from '../reducers/Wizard.reducer';

const initialState = {
    [STORE_SECTION_FORM]: {
        ccAttributes: {},
        SELECTED_MANAGEMENT_CLUSTER: '',
        SELECTED_CLUSTER_CLASS: '',
        AVAILABLE_CLUSTER_CLASSES: [],
    },
    [STORE_SECTION_UI]: {
        wcCcCategoryExpanded: {},
    },
};

const WcStore = createContext<{
    state: { [key: string]: any };
    dispatch: StoreDispatch;
}>({
    state: initialState,
    dispatch: () => null,
});

const WcProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode }) => {
    const [state, dispatch] = useReducer(wcReducer, initialState);

    return <WcStore.Provider value={{ state, dispatch }}>{children}</WcStore.Provider>;
};

export { WcStore, WcProvider };

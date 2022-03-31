import React, { createContext, Dispatch, ReactNode, useReducer } from 'react';
import mainReducer from '../reducers';

const initialState = {
    app: {
        appEnv: ''
    },
    ui: {
        navExpanded: false
    },
    data: {
        VCENTER_SERVER: 'abcd'
    }
};
const Store = createContext<{
    state: any,
    dispatch: Dispatch<any>
}>({
    state: initialState,
    dispatch: () => null
});


const AppProvider: React.FC<{ children: ReactNode}> = ({ children } : { children: ReactNode}) => {
    const [state, dispatch] = useReducer(mainReducer, initialState);

    return (
        <Store.Provider value={{ state, dispatch }}>
            {children}
        </Store.Provider>
    );
};

export { Store, AppProvider };
// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import dockerReducer from '../../../../state-management/reducers/Wizard.reducer';
import { DOCKER_DEFAULT_VALUES } from '../../../../shared/constants/defaults/docker.defaults';
import { StoreDispatch } from '../../../../shared/types/types';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';

const initialState = {
    [STORE_SECTION_FORM]: {
        ...DOCKER_DEFAULT_VALUES,
    },
};

const DockerStore = createContext<{
    dockerState: { [key: string]: any };
    dockerDispatch: StoreDispatch;
}>({
    dockerState: initialState,
    dockerDispatch: () => null,
});

const DockerProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode }) => {
    const [dockerState, dockerDispatch] = useReducer(dockerReducer, initialState);

    return <DockerStore.Provider value={{ dockerState, dockerDispatch }}>{children}</DockerStore.Provider>;
};

export { DockerStore, DockerProvider };

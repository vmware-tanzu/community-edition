// React imports
import React, { createContext, ReactNode, useReducer } from 'react';
// App imports
import { AWS_DEFAULT_VALUES } from '../../shared/constants/defaults/aws.defaults';
import awsReducer from '../reducers/Wizard.reducer';
import { STORE_SECTION_FORM } from '../reducers/Form.reducer';
import { StoreDispatch } from '../../shared/types/types';

const initialState = {
    [STORE_SECTION_FORM]: {
        // Auth - Credential Profile
        PROFILE: '',

        // Auth - Temporary Credentials
        SECRET_ACCESS_KEY: '',
        SESSION_TOKEN: '',
        ACCESS_KEY_ID: '',

        // Auth - General
        EC2_KEY_PAIR: '',

        // InstanceType
        NODE_PROFILE: '',

        // Region
        REGION: '',
        ...AWS_DEFAULT_VALUES,
    },
};

const AwsStore = createContext<{
    awsState: { [key: string]: any };
    awsDispatch: StoreDispatch;
}>({
    awsState: initialState,
    awsDispatch: () => null,
});

const AwsProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode }) => {
    const [awsState, awsDispatch] = useReducer(awsReducer, initialState);

    return <AwsStore.Provider value={{ awsState, awsDispatch }}>{children}</AwsStore.Provider>;
};

export { AwsStore, AwsProvider };

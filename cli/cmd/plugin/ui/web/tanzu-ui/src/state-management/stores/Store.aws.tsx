// React imports
import React, {
    createContext,
    ReactNode,
    useReducer,
} from 'react';

// App imports
import { StoreDispatch } from '../../shared/types/types';
import awsReducer from '../reducers/Aws.reducer';

const initialState = {
    data: {
        PROFILE: '',
        SECRET_ACCESS_KEY: '',
        SESSION_TOKEN: '',
        ACCESS_KEY_ID: '',
        REGION: '',
        EC2_KEY_PAIR: '',
        CLUSTER_NAME: ''
    },
};

const AwsStore = createContext<{
    awsState: { [key: string]: any };
    awsDispatch: StoreDispatch;
}>({
    awsState: initialState,
    awsDispatch: () => null,
});

const AwsProvider: React.FC<{ children: ReactNode }> = ({
    children,
}: {
    children: ReactNode;
}) => {
    const [awsState, awsDispatch] = useReducer(
        awsReducer,
        initialState
    );

    return (
        <AwsStore.Provider value={{ awsState, awsDispatch }}>
            {children}
        </AwsStore.Provider>
    );
};

export { AwsStore, AwsProvider };

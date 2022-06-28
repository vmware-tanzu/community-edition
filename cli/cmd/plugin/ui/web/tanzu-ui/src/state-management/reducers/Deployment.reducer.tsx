// App imports
import { Action } from '../../shared/types/types';
import { DEPLOYMENT_STATUS_CHANGED } from '../actions/Deployment.actions';
import { Deployments } from '../../shared/models/Deployments';
import { ReducerDescriptor } from '../../shared/utilities/Reducer.utils';

export const STORE_SECTION_DEPLOYMENT = 'deployment';

interface FormState {
    deployments: Deployments;
}

function deploymentReducer(state: FormState, action: Action) {
    const newState = { ...state, deployments: { ...action.payload } };
    console.log(`New deployment state: ${JSON.stringify(newState)}`);
    return newState;
}

export const deploymentReducerDescriptor = {
    name: 'deployment reducer',
    reducer: deploymentReducer,
    actionTypes: [DEPLOYMENT_STATUS_CHANGED],
    storeSection: STORE_SECTION_DEPLOYMENT,
} as ReducerDescriptor;

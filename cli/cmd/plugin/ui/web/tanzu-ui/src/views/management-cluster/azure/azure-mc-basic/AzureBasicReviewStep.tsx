// React imports
import React, { useContext } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';

// App imports
import { AzureStore } from '../store/Azure.store.mc';
import { AzureConfigDisplayConfig, AzureConfigDisplayDefaults } from './AzureBasicReview.config';
import { ConfigDisplay } from '../../../../shared/components/ConfigReview/ConfigDisplay';
import { NavRoutes } from '../../../../shared/constants/NavRoutes.constants';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';

export function AzureBasicReviewStep(props: Partial<StepProps>) {
    const { azureState } = useContext(AzureStore);
    const navigate = useNavigate();

    const navigateToDeploymentProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    const handleMCCreation = () => {
        props.deploy && props.deploy();
        navigateToDeploymentProgress();
    };
    return (
        <div cds-layout="p:md">
            <div cds-text="section" cds-layout="p-y:xs">
                Review the cluster configuration
            </div>
            <ConfigDisplay data={AzureConfigDisplayConfig} startsOpen={true} store={azureState[STORE_SECTION_FORM]} />
            <ConfigDisplay data={AzureConfigDisplayDefaults} store={azureState[STORE_SECTION_FORM]} />
            <CdsButton cds-layout="col:start-1" status="success" onClick={handleMCCreation}>
                <CdsIcon shape="cluster" size="sm"></CdsIcon>
                Create Management cluster
            </CdsButton>
        </div>
    );
}

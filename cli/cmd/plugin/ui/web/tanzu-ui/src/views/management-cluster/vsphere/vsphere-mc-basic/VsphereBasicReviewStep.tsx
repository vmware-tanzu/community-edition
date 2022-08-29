// React imports
import React, { useContext } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';

// App imports
import { configDisplayBasic, configDisplayDefaults } from './VsphereBasicReview.config';
import { ConfigDisplay } from '../../../../shared/components/ConfigReview/ConfigDisplay';
import { NavRoutes } from '../../../../shared/constants/NavRoutes.constants';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { VsphereStore } from '../Store.vsphere.mc';

export function VsphereBasicReviewStep(props: Partial<StepProps>) {
    const { vsphereState } = useContext(VsphereStore);
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
            <ConfigDisplay data={configDisplayBasic} startsOpen={true} store={vsphereState[STORE_SECTION_FORM]} />
            <ConfigDisplay data={configDisplayDefaults} store={vsphereState[STORE_SECTION_FORM]} />
            <CdsButton cds-layout="col:start-1" status="success" onClick={handleMCCreation}>
                <CdsIcon shape="cluster" size="sm"></CdsIcon>
                Create Management cluster
            </CdsButton>
        </div>
    );
}

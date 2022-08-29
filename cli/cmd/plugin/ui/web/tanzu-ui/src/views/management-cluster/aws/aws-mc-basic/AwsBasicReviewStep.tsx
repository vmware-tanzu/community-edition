// React imports
import React, { useContext } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';

// App imports
import { AwsStore } from '../store/Aws.store.mc';
import {
    AwsConfigDisplayBasicOneTimeCredentials,
    AwsConfigDisplayBasicProfileCredentials,
    AwsConfigDisplayDefaults,
} from './AwsBasicReview.config';
import { AWS_FIELDS, CREDENTIAL_TYPE } from './AwsManagementClusterBasic.constants';
import { ConfigDisplay, ConfigDisplayData } from '../../../../shared/components/ConfigReview/ConfigDisplay';
import { NavRoutes } from '../../../../shared/constants/NavRoutes.constants';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';

export function AwsBasicReviewStep(props: Partial<StepProps>) {
    const { awsState } = useContext(AwsStore);
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
            <ConfigDisplay
                data={createBasicDisplay(awsState[STORE_SECTION_FORM][AWS_FIELDS.CREDENTIAL_TYPE])}
                startsOpen={true}
                store={awsState[STORE_SECTION_FORM]}
            />
            <ConfigDisplay data={AwsConfigDisplayDefaults} store={awsState[STORE_SECTION_FORM]} />
            <CdsButton cds-layout="col:start-1" status="success" onClick={handleMCCreation}>
                <CdsIcon shape="cluster" size="sm"></CdsIcon>
                Create Management cluster
            </CdsButton>
        </div>
    );
}

function createBasicDisplay(credentialType: string): ConfigDisplayData {
    return credentialType === CREDENTIAL_TYPE.PROFILE ? AwsConfigDisplayBasicProfileCredentials : AwsConfigDisplayBasicOneTimeCredentials;
}

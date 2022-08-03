// React imports
import React, { useContext } from 'react';
import { useNavigate } from 'react-router-dom';
// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
// App imports
import { AWS_FIELDS, CREDENTIAL_TYPE } from './AwsManagementClusterBasic.constants';
import { AwsStore } from '../../../../state-management/stores/Store.aws';
import { ConfigDisplay, ConfigDisplayData } from '../../../../shared/components/ConfigReview/ConfigDisplay';
import { NavRoutes } from '../../../../shared/constants/NavRoutes.constants';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import {
    AwsConfigDisplayBasicOneTimeCredentials,
    AwsConfigDisplayBasicProfileCredentials,
    AwsConfigDisplayDefaults,
} from './AwsBasicReview.config';

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
        <>
            <div>Review the cluster configuration</div>
            <hr />
            <ConfigDisplay
                data={createBasicDisplay(awsState[STORE_SECTION_FORM][AWS_FIELDS.CREDENTIAL_TYPE])}
                startsOpen={true}
                store={awsState[STORE_SECTION_FORM]}
            />
            <hr />
            <ConfigDisplay data={AwsConfigDisplayDefaults} store={awsState[STORE_SECTION_FORM]} />
            <hr />
            <CdsButton cds-layout="col:start-1" status="success" onClick={handleMCCreation}>
                <CdsIcon shape="cluster" size="sm"></CdsIcon>
                Create Management cluster
            </CdsButton>
        </>
    );
}

function createBasicDisplay(credentialType: string): ConfigDisplayData {
    return credentialType === CREDENTIAL_TYPE.PROFILE ? AwsConfigDisplayBasicProfileCredentials : AwsConfigDisplayBasicOneTimeCredentials;
}

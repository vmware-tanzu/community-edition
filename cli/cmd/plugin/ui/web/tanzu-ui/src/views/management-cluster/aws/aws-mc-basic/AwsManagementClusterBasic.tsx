// React imports
import React from 'react';

// App imports
import { AwsBasicReviewStep } from './AwsBasicReviewStep';
import AwsClusterSettingsStep from '../AwsClusterSettingsStep';
import { AWS_MC_BASIC_TAB_NAMES } from './AwsManagementClusterBasic.constants';
import ManagementCredentials from '../aws-mc-common/management-credential-step/ManagementCredentials';
import useAwsDeployment from '../../../../shared/services/awsDeployment';
import Wizard from '../../../../shared/components/wizard/Wizard';

function AwsManagementClusterBasic() {
    const { deployOnAws } = useAwsDeployment();

    return (
        <Wizard tabNames={AWS_MC_BASIC_TAB_NAMES}>
            <ManagementCredentials />
            <AwsClusterSettingsStep />
            <AwsBasicReviewStep deploy={deployOnAws} />
        </Wizard>
    );
}

export default AwsManagementClusterBasic;

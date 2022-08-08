// React imports
import React from 'react';

// App imports
import { AWS_MC_BASIC_TAB_NAMES } from './AwsManagementClusterBasic.constants';
import { AwsBasicReviewStep } from './AwsBasicReviewStep';
import AwsClusterSettingsStep from '../AwsClusterSettingsStep';
import ManagementCredentials from '../aws-mc-common/management-credential-step/ManagementCredentials';
import useAwsDeployment from '../../../../shared/services/awsDeployment';
import Wizard from '../../../../shared/components/wizard/Wizard';
import { AwsDefaults } from '../aws-mc-common/default-service/AwsDefaults.service';

function AwsManagementClusterBasic() {
    const { deployOnAws } = useAwsDeployment();
    const awsDefaultService = AwsDefaults();

    return (
        <Wizard tabNames={AWS_MC_BASIC_TAB_NAMES}>
            <ManagementCredentials defaultService={awsDefaultService} />
            <AwsClusterSettingsStep />
            <AwsBasicReviewStep deploy={deployOnAws} />
        </Wizard>
    );
}

export default AwsManagementClusterBasic;

// React imports
import React from 'react';
// App imports
import useVSphereDeployment from '../../../../shared/services/vSphereDeployment';
import { VSPHERE_MC_BASIC_TAB_NAMES } from './VsphereManagementClusterBasic.constants';
import { VsphereBasicReviewStep } from './VsphereBasicReviewStep';
import { VsphereCredentialsStep } from '../vsphere-mc-common/VsphereCredentialsStep';
import { VsphereClusterSettingsStep } from '../vsphere-mc-common/VsphereClusterSettingsStep';
import { VsphereClusterResourcesStep } from '../vsphere-mc-common/VsphereClusterResourcesStep';
import { VsphereLoadBalancerStep } from '../vsphere-mc-common/VsphereLoadBalancerStep';
// App imports
import Wizard from '../../../../shared/components/wizard/Wizard';

function VsphereManagementClusterBasic() {
    const { deployOnVSphere } = useVSphereDeployment();

    return (
        <Wizard tabNames={VSPHERE_MC_BASIC_TAB_NAMES}>
            <VsphereCredentialsStep />
            <VsphereClusterSettingsStep />
            <VsphereLoadBalancerStep />
            <VsphereClusterResourcesStep />
            <VsphereBasicReviewStep deploy={deployOnVSphere} />
        </Wizard>
    );
}

export default VsphereManagementClusterBasic;

// React imports
import React from 'react';

import { VSPHERE_MC_BASIC_TAB_NAMES } from './VsphereManagementClusterBasic.constants';
import { VsphereClusterResourcesStep } from '../vsphere-mc-common/VsphereClusterResourcesStep';
import { VsphereClusterSettingsStep } from '../vsphere-mc-common/VsphereClusterSettingsStep';
import { VsphereCredentialsStep } from '../vsphere-mc-common/VsphereCredentialsStep';
import { VsphereLoadBalancerStep } from '../vsphere-mc-common/VsphereLoadBalancerStep';
// App imports
import Wizard from '../../../../shared/components/wizard/Wizard';
import useVSphereDeployment from '../../../../shared/services/vSphereDeployment';

function VsphereManagementClusterBasic() {
    const { deployOnVSphere } = useVSphereDeployment();

    return (
        <Wizard tabNames={VSPHERE_MC_BASIC_TAB_NAMES}>
            <VsphereCredentialsStep />
            <VsphereClusterSettingsStep />
            <VsphereLoadBalancerStep />
            <VsphereClusterResourcesStep deploy={deployOnVSphere} />
        </Wizard>
    );
}

export default VsphereManagementClusterBasic;

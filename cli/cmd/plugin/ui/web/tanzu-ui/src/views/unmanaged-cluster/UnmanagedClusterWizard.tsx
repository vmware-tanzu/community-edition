// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsIconButton, CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { arrowIcon, ClarityIcons } from '@cds/core/icon';

// App imports
import { UmcProvider } from '../../state-management/stores/Store.umc';
import UnmanagedClusterWizardBasic from './unmanaged-cluster-wizard-page/UnmanagedClusterWizardBasic';
import UnmanagedClusterWizardAdvanced from './unmanaged-cluster-wizard-page/UnmanagedClusterWizardAdvanced';
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';

ClarityIcons.addIcons(arrowIcon);

function UnmanagedClusterWizard() {
    const navigate = useNavigate();
    // temp local variable for testing advanced settings, will be refactored
    const setting_switch = true;

    return (
        <UmcProvider>
            <div className="aws-management-container" cds-layout="grid col:12">
                <div cds-layout="col:8">
                    <div cds-layout="col:12 p-b:lg">
                        <div cds-layout="horizontal align:vertical-center">
                            <CdsIconButton action="flat" status="primary" onClick={() => navigate(NavRoutes.UNMANAGED_CLUSTER_INVENTORY)}>
                                <CdsIcon shape="arrow" direction="left" size="lg"></CdsIcon>
                            </CdsIconButton>
                            <span cds-text="title">Create Unmanaged Cluster</span>
                        </div>
                    </div>
                    <div cds-layout="vertical align:stretch">
                        <div className="section-raised" cds-layout="horizontal align:vertical-center p:md">
                            <div>Simple configuration</div>
                            <CdsButton action="outline" cds-layout="align:right" size="sm">
                                Use Advanced Configuration
                            </CdsButton>
                        </div>
                        {setting_switch ? <UnmanagedClusterWizardBasic /> : <UnmanagedClusterWizardAdvanced />}
                    </div>
                </div>
                <div cds-layout="col:4" className="image-container">
                    <div className="mgmt-cluster-admins-img"></div>
                </div>
            </div>
        </UmcProvider>
    );
}

export default UnmanagedClusterWizard;

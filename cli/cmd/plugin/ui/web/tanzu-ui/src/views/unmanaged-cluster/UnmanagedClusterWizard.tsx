// React imports
import React, { useState } from 'react';
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
    const [useBasicSettings, setUseBasicSettings] = useState<boolean>(true);

    return (
        <UmcProvider>
            <div cds-layout="grid col:12">
                <div cds-layout="col:8">
                    <div cds-layout="col:12 p-b:lg">
                        <div cds-layout="horizontal align:vertical-center">
                            <CdsIconButton action="flat" status="primary" onClick={() => navigate(NavRoutes.UNMANAGED_CLUSTER_INVENTORY)}>
                                <CdsIcon shape="arrow" direction="left" size="lg"></CdsIcon>
                            </CdsIconButton>
                            <span cds-text="title">Create Unmanaged Cluster</span>
                        </div>
                    </div>
                    <div cds-layout="vertical align:stretch">{RenderBasicAdvanced()}</div>
                </div>
                <div cds-layout="col:4" className="image-container">
                    <div className="mgmt-cluster-admins-img"></div>
                </div>
            </div>
        </UmcProvider>
    );

    function RenderBasicAdvanced() {
        return useBasicSettings ? UnmanagedBasicSetting() : UnmanagedAdvancedSetting();
    }

    function toggleBasicAdvanced() {
        setUseBasicSettings(!useBasicSettings);
    }

    function UnmanagedBasicSetting() {
        return (
            <div cds-layout="vertical align:stretch">
                <div className="section-raised" cds-layout="horizontal align:vertical-center p:md">
                    <div>Basic configuration</div>
                    <CdsButton action="outline" cds-layout="align:right" size="sm" onClick={() => toggleBasicAdvanced()}>
                        Use Advanced Configuration
                    </CdsButton>
                </div>
                <UnmanagedClusterWizardBasic />
            </div>
        );
    }

    function UnmanagedAdvancedSetting() {
        return (
            <div cds-layout="vertical align:stretch">
                <div className="section-raised" cds-layout="horizontal align:vertical-center p:md">
                    <div>Advanced configuration</div>
                    <CdsButton action="outline" cds-layout="align:right" size="sm" onClick={() => setUseBasicSettings(!useBasicSettings)}>
                        Use Basic Configuration
                    </CdsButton>
                </div>
                <UnmanagedClusterWizardAdvanced />
            </div>
        );
    }
}

export default UnmanagedClusterWizard;

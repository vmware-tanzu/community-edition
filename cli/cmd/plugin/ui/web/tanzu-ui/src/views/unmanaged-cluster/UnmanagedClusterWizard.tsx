// React imports
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { arrowIcon, ClarityIcons } from '@cds/core/icon';
import { CdsIconButton, CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';

// App imports
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import { UmcProvider } from '../../state-management/stores/Store.umc';
import UnmanagedClusterBasic from './unmanaged-cluster-basic/UnmanagedClusterBasic';
import UnmanagedClusterAdvanced from './unmanaged-cluster-advanced/UnmanagedClusterAdvanced';

ClarityIcons.addIcons(arrowIcon);

function UnmanagedClusterWizard() {
    const navigate = useNavigate();

    // temp local variable for testing advanced settings, will be refactored
    const [useBasicSettings, setUseBasicSettings] = useState<boolean>(true);

    return (
        <UmcProvider>
            <div cds-layout="col:12">
                <div cds-layout="p-b:lg">
                    <div cds-layout="horizontal align:vertical-center">
                        <CdsIconButton action="flat" status="primary" onClick={() => navigate(NavRoutes.UNMANAGED_CLUSTER_INVENTORY)}>
                            <CdsIcon shape="arrow" direction="left" size="lg"></CdsIcon>
                        </CdsIconButton>
                        <CdsIcon cds-layout="m-r:sm" shape="computer" size="xl" className="icon-blue"></CdsIcon>
                        <span cds-text="title">Create an Unmanaged Cluster</span>
                    </div>
                </div>
                <div cds-layout="vertical align:stretch">{RenderBasicAdvanced()}</div>
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
                <UnmanagedClusterBasic />
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
                <UnmanagedClusterAdvanced />
            </div>
        );
    }
}

export default UnmanagedClusterWizard;

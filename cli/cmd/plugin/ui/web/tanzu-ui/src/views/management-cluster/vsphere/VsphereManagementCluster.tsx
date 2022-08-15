// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { CdsIconButton } from '@cds/react/button';

// App imports
import VsphereLogo from '../../../assets/vsphere.svg';
import VsphereManagementClusterBasic from './vsphere-mc-basic/VsphereManagementClusterBasic';
import { VsphereProvider } from './Store.vsphere.mc';
import { NavRoutes } from '../../../shared/constants/NavRoutes.constants';

function VsphereManagementCluster() {
    const navigate = useNavigate();

    return (
        <VsphereProvider>
            <div className="vsphere-management-container" cds-layout="grid col:12">
                <div cds-layout="col:12">
                    <div cds-layout="col:12 p-b:lg">
                        <span cds-text="title">
                            <CdsIconButton
                                cds-layout="p-t:md"
                                action="flat"
                                status="primary"
                                onClick={() => navigate(NavRoutes.MANAGEMENT_CLUSTER_SELECT_PROVIDER)}
                            >
                                <CdsIcon shape="arrow" direction="left" size="lg"></CdsIcon>
                            </CdsIconButton>
                            <img src={VsphereLogo} className="logo logo-42" cds-layout="m-r:md" alt="vsphere logo" />
                            Create a Management Cluster on vSphere
                        </span>
                    </div>

                    {/* Disable Basic/Advanced banner until advanced settings available */}
                    {/* <ConfigBanner /> */}
                    <VsphereManagementClusterBasic />
                </div>
            </div>
        </VsphereProvider>
    );
}

export default VsphereManagementCluster;

// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { CdsIconButton } from '@cds/react/button';

// App imports
import { AzureProvider } from '../../../state-management/stores/Azure.store';
import AzureLogo from '../../../assets/azure.svg';
import AzureManagementClusterBasic from './azure-mc-basic/AzureManagementClusterBasic';
import './AzureManagementCluster.scss';
import { NavRoutes } from '../../../shared/constants/NavRoutes.constants';

function AzureManagementCluster() {
    const navigate = useNavigate();

    return (
        <AzureProvider>
            <div cds-layout="grid col:12">
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
                            <img src={AzureLogo} className="logo logo-42" cds-layout="m-r:md" alt="azure logo" />
                            Create a Management Cluster on Azure
                        </span>
                    </div>

                    {/* Disable Basic/Advanced banner until advanced settings available */}
                    {/* <ConfigBanner /> */}
                    <AzureManagementClusterBasic />
                </div>
            </div>
        </AzureProvider>
    );
}

export default AzureManagementCluster;

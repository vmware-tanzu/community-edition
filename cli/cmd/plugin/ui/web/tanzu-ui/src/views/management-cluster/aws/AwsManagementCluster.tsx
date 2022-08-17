// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { CdsIconButton } from '@cds/react/button';

// App imports
import AwsLogo from '../../../assets/aws.svg';
import AwsManagementClusterBasic from './aws-mc-basic/AwsManagementClusterBasic';
import { AwsProvider } from './store/Aws.store.mc';
import { NavRoutes } from '../../../shared/constants/NavRoutes.constants';
import './AwsManagementCluster.scss';

function AwsManagementCluster() {
    const navigate = useNavigate();

    return (
        <AwsProvider>
            <div className="aws-management-container" cds-layout="col:12">
                <div cds-layout="p-b:lg">
                    <span cds-text="title">
                        <CdsIconButton
                            cds-layout="p-t:md"
                            action="flat"
                            status="primary"
                            onClick={() => navigate(NavRoutes.MANAGEMENT_CLUSTER_SELECT_PROVIDER)}
                        >
                            <CdsIcon shape="arrow" direction="left" size="lg"></CdsIcon>
                        </CdsIconButton>
                        <img src={AwsLogo} className="logo logo-42" cds-layout="m-r:md" alt="aws logo" />
                        Create a Management Cluster on Amazon Web Services
                    </span>
                </div>

                {/* Disable Basic/Advanced banner until advanced settings available */}
                {/* <ConfigBanner /> */}
                <AwsManagementClusterBasic />
            </div>
        </AwsProvider>
    );
}

export default AwsManagementCluster;

// React imports
import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';

// App imports
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import { UnmanagedCluster } from '../../swagger-api';
import { UnmanagedService } from '../../swagger-api/services/UnmanagedService';
import UnmanagedClusterCard from './UnmanagedClusterCard/UnmanagedClusterCard';
import './UnmanagedClusterInventory.scss';

function UnmanagedClusterInventory() {
    const [unmanagedClusters, setUnmanagedClusters] = useState<UnmanagedCluster[]>([]);
    const navigate = useNavigate();

    useEffect(() => {
        // fetch unmanaged clusters
        UnmanagedService.getUnmanagedClusters().then((data) => setUnmanagedClusters(data));
    }, []);

    return (
        <div className="management-cluster-landing-container" cds-layout="grid vertical col:12 gap:lg align:fill">
            <div cds-layout="grid horizontal col:12">
                <div cds-layout="vertical gap:md gap@md:lg col@sm:8 col:8">
                    <div cds-text="title" cds-layout="horizontal align:vertical-center">
                        <CdsIcon cds-layout="m-r:sm" shape="computer" size="lg"></CdsIcon>
                        Unmanaged Clusters
                    </div>
                    <div cds-text="body">
                        Create a single node, local workstation cluster suitable for a development/test environment. It requires minimal
                        local resources and is fast to deploy. It provides support for running multiple clusters. The default Tanzu
                        Community Edition package repository is automatically installed when you deploy an unmanaged cluster.
                    </div>
                    <div cds-layout="vertical gap:md gap@md:lg col@sm:12 col:12">
                        <CdsButton
                            className="cluster-action-btn"
                            status="primary"
                            onClick={() => navigate(NavRoutes.UNMANAGED_CLUSTER_WIZARD)}
                        >
                            <CdsIcon shape="block"></CdsIcon>
                            Create an Unmanaged Cluster
                        </CdsButton>
                    </div>
                    <div cds-layout="vertical gap:lg col:6">
                        <div cds-layout="vertical gap:md gap@md:lg col@sm:12 col:12">
                            <CdsAlertGroup status="success">
                                <CdsAlert>This is an alert with a status</CdsAlert>
                            </CdsAlertGroup>
                        </div>
                        {unmanagedClusters.map((data, index) => {
                            return <UnmanagedClusterCard key={index} name={data.name} provider={data.provider} status={data.status} />;
                        })}
                    </div>
                </div>
                <div cds-layout="col@sm:4 col:4 container:fill"></div>
            </div>
        </div>
    );
}

export default UnmanagedClusterInventory;

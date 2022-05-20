// React imports
import React from 'react';
// Library imports
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';
// App imports
import { ManagementCluster } from '../../swagger-api';
import { ProviderData, retrieveProviderInfo } from '../../shared/services/Provider.service';

function ManagementClusterInfoBanner(managementCluster: ManagementCluster) {
    if (!managementCluster) {
        return <></>;
    }
    const provider: ProviderData = retrieveProviderInfo(managementCluster.provider || 'unknown');
    return (
        <CdsAlertGroup
            type="banner"
            status="neutral"
            aria-label={`This workload cluster will be provisioned on ${managementCluster.provider} using ${managementCluster.name}`}
        >
            <CdsAlert closable>
                This workload cluster will be provisioned on {managementCluster.provider} using <b>{managementCluster.name}</b> &nbsp;
                <img src={provider.logo} className="logo logo-26" alt={`${provider.name} logo`} />
            </CdsAlert>
        </CdsAlertGroup>
    );
}

export default ManagementClusterInfoBanner;

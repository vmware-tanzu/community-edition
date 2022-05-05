// React imports
import React from 'react';
// Library imports
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';
// App imports
import { ManagementCluster } from '../../shared/models/ManagementCluster';

function ManagementClusterInfoBanner(managementCluster: ManagementCluster) {
    if (!managementCluster) {
        return <></>
    }
    return <CdsAlertGroup
        type="banner"
        status="success"
        aria-label={`This workload cluster will be provisioned on ${managementCluster.provider} using ${managementCluster.name}`}
    >
        <CdsAlert closable>
            This workload cluster will be provisioned on {managementCluster.provider} using <b>{managementCluster.name}</b>
        </CdsAlert>
    </CdsAlertGroup>;
}

export default ManagementClusterInfoBanner;

// React imports
import React from 'react';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, cloudIcon } from '@cds/core/icon';

// App imports

ClarityIcons.addIcons(cloudIcon);

const ManagementClusterInventory: React.FC = () => {
    return (
        <div className="management-cluster-landing-container" cds-layout="vertical gap:md col@sm:12 grid">
            <div cds-layout="vertical col:8 gap:lg">
                <div cds-text="title">
                    <CdsIcon cds-layout="m-r:sm" shape="cloud" size="xl" className="icon-blue"></CdsIcon>
                    Management Cluster Inventory
                </div>
                <div cds-text="subsection" cds-layout="p-y:md">
                    Placeholder page for management cluster inventory
                </div>
            </div>
            <div cds-layout="col:4" className="mgmt-cluster-admins-img"></div>
        </div>
    );
};

export default ManagementClusterInventory;

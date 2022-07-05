// React imports
import React from 'react';

// Library imports
import { CdsIcon } from '@cds/react/icon';

function WorkloadClusterInventory() {
    return (
        <div cds-layout="grid vertical col:12 gap:xl align:fill">
            <div cds-layout="grid horizontal col:12">
                <div cds-layout="vertical col:8 gap:lg">
                    <div cds-text="title">
                        <CdsIcon cds-layout="m-r:sm" shape="nodes" size="xl" className="icon-blue"></CdsIcon>
                        Workload Cluster Inventory
                    </div>
                    <div cds-text="subsection" cds-layout="p-y:md">
                        Placeholder page for workload cluster inventory
                    </div>
                </div>
            </div>
        </div>
    );
}

export default WorkloadClusterInventory;

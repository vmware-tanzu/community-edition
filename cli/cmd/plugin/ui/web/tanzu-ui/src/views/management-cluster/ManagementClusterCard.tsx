// React imports
import React from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsDivider } from '@cds/react/divider';
import { CdsIcon } from '@cds/react/icon';

// App imports
import { ManagementCluster } from '../../swagger-api';

interface ManagementClusterProps extends ManagementCluster {
    confirmDeleteCallback: (arg?: string) => void;
}

function ManagementClusterCard(props: ManagementClusterProps) {
    const { name, path, context, confirmDeleteCallback } = props;

    return (
        <div className="section-raised" cds-layout="grid cols:12 wrap:none" data-testid="management-cluster-card">
            <div cds-layout="vertical">
                <div cds-layout="vertical" cds-text="subsection">
                    <div cds-layout="horizontal gap:sm align:vertical-center p:sm">
                        <CdsIcon shape="blocks-group" size="md" className="icon-blue"></CdsIcon>
                        <div cds-text="section" cds-layout="align:stretch">
                            {name}
                        </div>
                        <CdsButton
                            action="flat-inline"
                            status="danger"
                            onClick={() => {
                                confirmDeleteCallback(name);
                            }}
                        >
                            Delete
                        </CdsButton>
                    </div>
                    <CdsDivider></CdsDivider>
                    <div cds-layout="horizontal gap:md p:sm">
                        <div cds-layout="vertical">
                            <label cds-text="p4" cds-layout="m-b:sm">
                                Path
                            </label>
                            <div cds-layout="horizontal">
                                <div cds-text="body">{path}</div>
                            </div>
                        </div>
                        <CdsDivider orientation="vertical"></CdsDivider>
                        <div cds-layout="vertical m-r:xs">
                            <label cds-text="p4" cds-layout="m-b:sm">
                                Context
                            </label>
                            <div cds-layout="horizontal">
                                <div cds-text="body">{context}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default ManagementClusterCard;

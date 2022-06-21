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
                <div cds-layout="horizontal gap:md align:fill align:vertical-center p:md" cds-text="subsection">
                    <div cds-layout="horizontal">
                        <div cds-layout="horizontal gap:sm align:vertical-center p-y:sm">
                            <CdsIcon cds-layout="m-r:sm" shape="cluster" size="lg" className="icon-blue"></CdsIcon>
                            <div cds-text="section">{name}</div>
                        </div>
                        <CdsDivider orientation="vertical" cds-layout="align:right"></CdsDivider>
                    </div>
                    <div cds-layout="horizontal">
                        <div cds-layout="vertical align:left m-r:xs">
                            <label className="card-content-label">Path</label>
                            <div cds-layout="horizontal">
                                <div>{path}</div>
                            </div>
                        </div>
                        <CdsDivider orientation="vertical" cds-layout="align:right"></CdsDivider>
                    </div>
                    <div cds-layout="horizontal">
                        <div cds-layout="vertical m-r:xs">
                            <label className="card-content-label">Context</label>
                            <div cds-layout="horizontal">
                                <div>{context}</div>
                            </div>
                        </div>
                    </div>
                </div>
                <CdsDivider></CdsDivider>
                <div cds-layout="vertical gap:md align:vertical-center p:md">
                    <div cds-layout="horizontal gap:xs align:vertical-center align:right p-r:lg">
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
                </div>
            </div>
        </div>
    );
}

export default ManagementClusterCard;

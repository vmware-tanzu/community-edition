// React imports
import React from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsDivider } from '@cds/react/divider';
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, disconnectIcon, unknownStatusIcon } from '@cds/core/icon';

// App imports
import { UnmanagedCluster } from '../../../../swagger-api';
import './UnmanagedClusterCard.scss';

ClarityIcons.addIcons(disconnectIcon, unknownStatusIcon);

/* eslint-disable no-unused-vars */
enum UnmanagedClusterStatus {
    RUNNING = 'RUNNING',
    STOPPED = 'STOPPED',
    UNKNOWN = 'UNKNOWN',
}

interface UnmanagedClusterProps extends UnmanagedCluster {
    confirmDeleteCallback: (arg?: string) => void;
}

function UnmanagedClusterCard(props: UnmanagedClusterProps) {
    const { name, provider, status, confirmDeleteCallback } = props;

    return (
        <div
            className={'section-raised status-' + status?.toLowerCase()}
            cds-layout="grid cols:12 wrap:none"
            data-testid="unmanaged-cluster-card"
        >
            <div cds-layout="vertical">
                <div cds-layout="vertical" cds-text="subsection">
                    <div cds-layout="horizontal gap:sm align:vertical-center p:sm">
                        <CdsIcon shape="block" size="md" className="icon-blue"></CdsIcon>
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
                        <div cds-layout="vertical m-r:xs">
                            <label cds-text="p4" cds-layout="m-b:sm">
                                Provider
                            </label>
                            <div cds-layout="horizontal">
                                <div cds-text="body">{provider}</div>
                            </div>
                        </div>
                        <CdsDivider orientation="vertical"></CdsDivider>
                        <div cds-layout="vertical m-r:xs">
                            <label cds-text="p4" cds-layout="m-b:sm">
                                Status
                            </label>
                            <div cds-layout="horizontal align:vertical-center">
                                {status?.toUpperCase() === UnmanagedClusterStatus.RUNNING ? (
                                    <CdsIcon cds-layout="m-r:xs" shape="check" size="sm" status="success"></CdsIcon>
                                ) : status?.toUpperCase() === UnmanagedClusterStatus.STOPPED ? (
                                    <CdsIcon cds-layout="m-r:xs" shape="exclamation-circle" solid size="sm" status="danger"></CdsIcon>
                                ) : (
                                    <CdsIcon cds-layout="m-r:xs" shape="unknown-status" size="sm"></CdsIcon>
                                )}
                                <div cds-text="body">{status}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default UnmanagedClusterCard;

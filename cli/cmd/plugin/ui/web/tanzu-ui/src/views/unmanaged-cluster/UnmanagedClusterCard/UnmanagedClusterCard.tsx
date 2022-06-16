// React imports
import React from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { CdsDivider } from '@cds/react/divider';
import { ClarityIcons, disconnectIcon, unknownStatusIcon } from '@cds/core/icon';

// App imports
import './UnmanagedClusterCard.scss';
import { UnmanagedCluster } from '../../../swagger-api';

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
        <div className="section-raised" cds-layout="grid cols:12 wrap:none" data-testid="unmanaged-cluster-card">
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
                            <label>Provider</label>
                            <div cds-layout="horizontal">
                                <div>{provider}</div>
                            </div>
                        </div>
                        <CdsDivider orientation="vertical" cds-layout="align:right"></CdsDivider>
                    </div>
                    <div cds-layout="horizontal">
                        <div cds-layout="vertical m-r:xs">
                            <label>Status</label>
                            <div cds-layout="horizontal">
                                {status?.toUpperCase() === UnmanagedClusterStatus.RUNNING ? (
                                    <CdsIcon cds-layout="m-r:sm" shape="connect" size="md" status="success"></CdsIcon>
                                ) : status?.toUpperCase() === UnmanagedClusterStatus.STOPPED ? (
                                    <CdsIcon cds-layout="m-r:sm" shape="disconnect" size="md" status="danger"></CdsIcon>
                                ) : (
                                    <CdsIcon cds-layout="m-r:sm" shape="unknown-status" size="md" status="warning"></CdsIcon>
                                )}
                                <div>{status}</div>
                            </div>
                        </div>
                    </div>
                </div>
                <CdsDivider></CdsDivider>
                <div cds-layout="vertical gap:md align:vertical-center p:md">
                    <div cds-layout="horizontal gap:xs align:vertical-center align:left">
                        <CdsButton action="flat" size="sm" cds-layout="m-x:md">
                            Access This Cluster
                        </CdsButton>
                        <CdsButton
                            action="flat"
                            size="sm"
                            cds-layout="m-x:md"
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

export default UnmanagedClusterCard;

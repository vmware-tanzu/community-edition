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

enum UnmanagedClusterStatus {
    RUNNING = 'RUNNING',
    STOPPED = 'STOPPED',
    UNKNOWN = 'UNKNOWN',
}

function UnmanagedClusterCard(props: UnmanagedCluster) {
    const { name, provider, status } = props;
    return (
        <div className="section-raised" cds-layout="grid cols:12 wrap:none">
            <div cds-layout="vertical">
                <div cds-layout="vertical" cds-text="subsection">
                    <div cds-layout="horizontal">
                        <div cds-layout="horizontal gap:sm align:vertical-center p:sm">
                            <CdsIcon shape="block" size="md" className="icon-blue"></CdsIcon>
                            <div cds-text="section">{name}</div>
                        </div>
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
                <CdsDivider></CdsDivider>
                <div cds-layout="vertical gap:md p:md">
                    <div cds-layout="horizontal gap:lg">
                        <CdsButton action="flat-inline" size="md">
                            Access This Cluster
                        </CdsButton>
                        <CdsButton action="flat-inline" size="md" status="danger">
                            Delete
                        </CdsButton>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default UnmanagedClusterCard;

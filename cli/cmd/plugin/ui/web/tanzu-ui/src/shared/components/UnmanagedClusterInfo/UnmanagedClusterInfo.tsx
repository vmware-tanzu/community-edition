// React imports
import React from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { CdsDivider } from '@cds/react/divider';

// App imports
import './UnmanagedClusterInfo.scss';

function UnmanagedClusterInfo(props: any) {
    return (
        <div className="section-raised" cds-layout="grid cols:12 wrap:none">
            <div cds-layout="vertical">
                <div cds-layout="horizontal gap:md align:fill align:vertical-center p:md" cds-text="subsection">
                    <div cds-layout="horizontal">
                        <div cds-layout="horizontal gap:sm align:vertical-center p-y:sm">
                            <CdsIcon cds-layout="m-r:sm" shape="cluster" size="lg" className="icon-blue"></CdsIcon>
                            <div cds-text="section">{props.name}</div>
                        </div>
                        <CdsDivider orientation="vertical" cds-layout="align:right"></CdsDivider>
                    </div>

                    <div cds-layout="horizontal">
                        <div cds-layout="vertical align:left m-r:xs">
                            <label>Provider</label>
                            <div cds-layout="horizontal">
                                <div>{props.provider}</div>
                            </div>
                        </div>
                        <CdsDivider orientation="vertical" cds-layout="align:right"></CdsDivider>
                    </div>
                    <div cds-layout="horizontal">
                        <div cds-layout="vertical m-r:xs">
                            <label>Status</label>
                            <div cds-layout="horizontal">
                                {props.status == 'Running' ? (
                                    <CdsIcon cds-layout="m-r:sm" shape="connect" size="md" status="success"></CdsIcon>
                                ) : (
                                    <CdsIcon cds-layout="m-r:sm" shape="disconnect" size="md" status="danger"></CdsIcon>
                                )}
                                <div>{props.status}</div>
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
                        <CdsButton action="flat" size="sm" cds-layout="m-x:md" status="danger">
                            Delete
                        </CdsButton>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default UnmanagedClusterInfo;

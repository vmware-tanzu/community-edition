// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, blockIcon, computerIcon } from '@cds/core/icon';

// App imports
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import './GettingStarted.scss';

ClarityIcons.addIcons(blockIcon, computerIcon);

const GettingStarted: React.FC = () => {
    const navigate = useNavigate();
    return (
        <>
            <div cds-layout="grid vertical col:12 gap:xl align:fill">
                <div cds-layout="grid horizontal col:12">
                    <div cds-layout="vertical gap:md gap@md:lg col@sm:8 col:8">
                        <div cds-text="title">
                            <CdsIcon cds-layout="m-r:sm" shape="block" size="xl" className="icon-blue"></CdsIcon>
                            Getting Started
                        </div>
                        <div cds-text="subsection">
                            Creating a full-featured, Kubernetes implementation suitable for a development or production
                            environment starts with a Management Cluster on a cloud provider. From the Management Cluster, you
                            will be able to create Workload Clusters.
                        </div>
                    </div>
                    <div cds-layout="col@sm:4 col:4 container:fill">
                        <CdsButton action="flat" onClick={() => {
                            window.open('http://tanzucommunityedition.io', '_blank');
                        }}>Learn more about Tanzu&apos;s architecture</CdsButton>
                    </div>
                </div>

                <div cds-layout="grid horizontal col:12">
                    <div cds-layout="vertical gap:lg gap@md:lg col@sm:6 col:6" className="mgmt-cluster-intro-container">
                        <div cds-text="title" className="text-blue">
                            Management Cluster
                        </div>
                        <div cds-layout="grid cols:12 gap:lg gap@md:lg">
                            <div cds-layout="grid cols:6 gap:lg gap@md:lg" cds-text="body">
                                <div cds-text="body">
                                    Creating a full-featured, Kubernetes implementation suitable for a development or production
                                    environment starts with a Management Cluster on a cloud provider. From the Management Cluster, you
                                    will be able to create Workload Clusters.
                                </div>
                            </div>
                            <div cds-layout="grid cols:12">
                                <CdsButton
                                    className="cluster-action-btn"
                                    status="primary"
                                    onClick={()=> navigate(NavRoutes.MANAGEMENT_CLUSTER_LANDING)}>
                                    Create a Management Cluster
                                </CdsButton>
                            </div>
                        </div>
                    </div>
                    <div cds-layout="vertical gap:lg gap@md:lg col@sm:6 col:6" className="wl-cluster-intro-container">
                        <div cds-text="title" className="text-blue">
                            Workload Cluster
                        </div>
                        <div cds-layout="grid cols:12 gap:lg gap@md:lg">
                            <div cds-layout="grid cols:6 gap:lg gap@md:lg">
                                <div cds-text="body">
                                    Creating a full-featured, Kubernetes implementation suitable for a development or production
                                    environment starts with a Management Cluster on a cloud provider. From the Management Cluster, you
                                    will be able to create Workload Clusters.
                                </div>
                            </div>
                            <div cds-layout="grid cols:12">
                                <CdsButton
                                    className="cluster-action-btn"
                                    status="neutral"
                                    onClick={()=> navigate(NavRoutes.WORKLOAD_CLUSTER_WIZARD)}>
                                    Create a Workload Cluster
                                </CdsButton>
                            </div>
                        </div>
                    </div>
                </div>

                <div cds-layout="grid vertical col:12 ">
                    <div cds-layout="vertical gap:md gap@md:lg col@sm:12 align:fill">
                        <div cds-text="title">
                            Need a local cluster for experimentation and development?
                        </div>
                        <div cds-layout="grid cols:12">
                            <div cds-text="subsection">
                                Unmanaged Clusters offer Tanzu environments for development and experimentation.<br/>
                                By default, they run locally with Tanzu components installed on top.
                            </div>
                        </div>
                        <div cds-layout="grid cols:12">
                            <CdsButton
                                className="cluster-action-btn"
                                action="outline"
                                onClick={()=> navigate(NavRoutes.UNMANAGED_CLUSTER_LANDING)}>
                                <CdsIcon shape="computer"></CdsIcon>
                                Create an Unmanaged Cluster
                            </CdsButton>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
};

export default GettingStarted;

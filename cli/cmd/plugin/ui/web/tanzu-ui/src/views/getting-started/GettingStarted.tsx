// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, blockIcon, computerIcon } from '@cds/core/icon';

// App imports
import './GettingStarted.scss';
import { AppFeature, featureAvailable } from '../../shared/services/AppConfiguration.service';
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';

ClarityIcons.addIcons(blockIcon, computerIcon);

function GettingStarted() {
    const navigate = useNavigate();
    const workloadClusterSupport = featureAvailable(AppFeature.WORKLOAD_CLUSTER_SUPPORT);
    const unmanagedClusterSupport = featureAvailable(AppFeature.UNMANAGED_CLUSTER_SUPPORT);
    return (
        <>
            <div cds-layout="grid vertical col:12 gap:lg align:fill">
                <div cds-layout="grid horizontal col:12">
                    <div cds-layout="vertical gap:md gap@md:lg col@sm:8 col:8">
                        <div cds-text="title">
                            <CdsIcon cds-layout="m-r:sm" shape="compass" size="xl" className="icon-aqua"></CdsIcon>
                            Getting started
                        </div>
                        <div cds-text="subsection">
                            <p>
                                Creating a full-featured, Kubernetes implementation suitable for a development or production environment
                                starts with a Management Cluster on a cloud provider.
                            </p>
                            <p>From the Management Cluster, you will be able to create Workload Clusters.</p>
                        </div>
                    </div>
                    <div cds-layout="col@sm:4 col:4 container:fill">
                        <CdsButton
                            action="flat"
                            onClick={() => {
                                window.open('http://tanzucommunityedition.io', '_blank');
                            }}
                        >
                            Learn more about Tanzu&apos;s architecture
                        </CdsButton>
                    </div>
                </div>

                <div cds-layout="grid horizontal col:12 gap:lg h:">
                    <div cds-layout="vertical gap:lg gap@md:lg col@sm:6 col:6 p:lg" className="section-raised mgmt-cluster-intro-container">
                        <div cds-text="title" className="text-aqua">
                            Management Cluster
                        </div>
                        <div cds-layout="grid cols:12 gap:lg gap@md:lg">
                            <div cds-layout="grid cols:7 gap:lg gap@md:lg" cds-text="body">
                                <div cds-text="body">
                                    <p cds-layout="m-t:none">
                                        Create a Management Cluster on your prefered cloud provider through a guided series of steps.
                                    </p>
                                    <p cds-layout="m-b:none">This cluster will manage new clusters you create for your workloads.</p>
                                </div>
                            </div>
                        </div>
                        <CdsButton
                            className="cluster-action-btn"
                            status="primary"
                            onClick={() => navigate(NavRoutes.MANAGEMENT_CLUSTER_SELECT_PROVIDER)}
                        >
                            Create a Management Cluster
                        </CdsButton>
                    </div>
                    {workloadClusterSupport && (
                        <div
                            cds-layout="vertical gap:lg gap@md:lg col@sm:6 col:6 p:lg"
                            className="section-raised wl-cluster-intro-container"
                        >
                            <div cds-text="title" className="text-green">
                                Workload Cluster
                            </div>
                            <div cds-layout="grid cols:12 gap:lg gap@md:lg align:stretch">
                                <div cds-layout="grid cols:7 gap:lg gap@md:lg">
                                    <div cds-text="body">
                                        <p cds-layout="m-t:none">Create a Workload Cluster from a Management Cluster.</p>
                                        <p cds-layout="m-b:none">Your workloads and Tanzu packages will be deployed to this cluster.</p>
                                    </div>
                                </div>
                            </div>
                            <CdsButton
                                className="cluster-action-btn"
                                action="outline"
                                disabled
                                status="neutral"
                                onClick={() => navigate(NavRoutes.WORKLOAD_CLUSTER_WIZARD)}
                            >
                                Create a Workload Cluster
                            </CdsButton>
                        </div>
                    )}
                </div>

                {unmanagedClusterSupport && (
                    <div cds-layout="grid vertical col:12 m-t:lg">
                        <div cds-layout="vertical gap:md gap@md:lg col@sm:12 align:fill" className="um-cluster-intro-container">
                            <div cds-text="title">Need a local cluster for experimentation and development?</div>
                            <div cds-layout="grid cols:4">
                                <p cds-text="body">
                                    Unmanaged Clusters offer Tanzu environments for development and experimentation. By default, they run
                                    locally with Tanzu components installed on top.
                                </p>
                            </div>
                            <CdsButton
                                className="cluster-action-btn"
                                action="outline"
                                onClick={() => navigate(NavRoutes.UNMANAGED_CLUSTER_INVENTORY)}
                            >
                                <CdsIcon shape="computer"></CdsIcon>
                                Create an Unmanaged Cluster
                            </CdsButton>
                        </div>
                    </div>
                )}
            </div>
        </>
    );
}

export default GettingStarted;

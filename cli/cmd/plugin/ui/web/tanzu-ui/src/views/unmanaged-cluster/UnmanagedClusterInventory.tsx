// React imports
import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';
import { CdsModal, CdsModalHeader, CdsModalActions, CdsModalContent } from '@cds/react/modal';

// App imports
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import { UnmanagedCluster } from '../../swagger-api';
import { UnmanagedService } from '../../swagger-api/services/UnmanagedService';
import UnmanagedClusterCard from './UnmanagedClusterCard/UnmanagedClusterCard';
import './UnmanagedClusterInventory.scss';

function UnmanagedClusterInventory() {
    const [unmanagedClusters, setUnmanagedClusters] = useState<UnmanagedCluster[]>([]);
    const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
    const [clusterNameForDeletion, setClusterNameForDeletion] = useState<string>('');
    const navigate = useNavigate();

    function retrieveUnmanagedClusters() {
        UnmanagedService.getUnmanagedClusters().then(
            (data) => {
                setUnmanagedClusters(data);
            },
            (err) => {
                console.warn(`unmanaged clusters get api failed with error: ${err}`);
            }
        );
    }

    useEffect(() => {
        retrieveUnmanagedClusters();
    }, []);

    function showConfirmDeleteModal(clusterName?: string) {
        if (!clusterName) {
            console.warn(`toggleConfirmDeleteModal expected a cluster name string; instead got ${clusterName}`);
            return;
        }
        setClusterNameForDeletion(clusterName);
        setShowDeleteModal(!showDeleteModal);
    }

    function hasUnmanagedClusters() {
        return unmanagedClusters.length ? true : false;
    }

    function deleteUnmanagedCluster(clusterName: string) {
        setShowDeleteModal(!showDeleteModal);
        UnmanagedService.deleteUnmanagedCluster(clusterName).then(
            () => {
                console.log(`unmanaged cluster ${clusterName} has been deleted`);
                retrieveUnmanagedClusters();
            },
            (err) => {
                console.log(`unmanaged cluster delete api failed with error: ${err}`);
            }
        );
        // TODO: show alert confirming deletion of mgmt cluster or failure - requires shared alert component
    }

    return (
        <div className="unmanaged-cluster-landing-container" cds-layout="grid vertical col:12 gap:lg align:fill">
            <div cds-layout="grid horizontal col:12">
                <div cds-layout="vertical gap:md gap@md:lg col@sm:8 col:8">
                    {Header()}
                    <CdsButton className="cluster-action-btn" status="primary" onClick={() => navigate(NavRoutes.UNMANAGED_CLUSTER_WIZARD)}>
                        <CdsIcon shape="cluster"></CdsIcon>
                        create unmanaged cluster
                    </CdsButton>
                    <CdsAlertGroup status="success">
                        <CdsAlert>This is an alert with a status</CdsAlert>
                    </CdsAlertGroup>
                    {hasUnmanagedClusters() ? UnmanagedClustersSection() : NoUnmanagedClustersSection()}
                </div>
                <div cds-layout="col@sm:4 col:4 container:fill"></div>
                {renderConfirmDeleteModal()}
            </div>
        </div>
    );

    function Header() {
        return (
            <div cds-layout="vertical gap:lg">
                <div cds-text="title">
                    <CdsIcon cds-layout="m-r:sm" shape="computer" size="xl" className="icon-blue"></CdsIcon>
                    Unmanaged Cluster
                </div>
                <div cds-text="subsection">
                    Create a single node, local workstation cluster suitable for a development/test environment. It requires minimal local
                    resources and is fast to deploy. It provides support for running multiple clusters. The default Tanzu Community Edition
                    package repository is automatically installed when you deploy an unmanaged cluster.
                </div>
            </div>
        );
    }

    function LearnMore() {
        return (
            <div cds-layout="vertical">
                <CdsButton
                    action="flat"
                    onClick={() => {
                        window.open('http://tanzucommunityedition.io', '_blank');
                    }}
                >
                    Learn more about Tanzu&apos;s architecture
                </CdsButton>
            </div>
        );
    }

    // Returns modal window HTML markup if showDeleteModal state variable is set to true.
    function renderConfirmDeleteModal() {
        if (!showDeleteModal) {
            return;
        }

        return (
            <>
                <CdsModal
                    aria-labelledby="default-modal-title"
                    id="confirm-delete-modal"
                    data-testid="confirm-delete-cluster-modal"
                    onCloseChange={() => setShowDeleteModal(false)}
                >
                    <CdsModalHeader>
                        <h3 cds-text="section" cds-first-focus="true" id="confirm-delete-modal-title">
                            Delete cluster {clusterNameForDeletion}
                        </h3>
                    </CdsModalHeader>
                    <CdsModalContent>
                        <p cds-text="body">
                            Deleting {clusterNameForDeletion} stops this Unmanaged Cluster and removes it from the provider you created it
                            on.
                        </p>
                    </CdsModalContent>
                    <CdsModalActions>
                        <CdsButton action="outline" onClick={() => setShowDeleteModal(false)}>
                            Cancel
                        </CdsButton>
                        <CdsButton
                            status="danger"
                            onClick={() => deleteUnmanagedCluster(clusterNameForDeletion)}
                            data-testid="delete-cluster-btn"
                        >
                            Delete
                        </CdsButton>
                    </CdsModalActions>
                </CdsModal>
            </>
        );
    }

    // Returns view to be rendered when Unamanged clusters are present.
    function UnmanagedClustersSection() {
        return (
            <>
                {unmanagedClusters.map((cluster: UnmanagedCluster) => {
                    return (
                        <UnmanagedClusterCard
                            key={cluster.name}
                            name={cluster.name}
                            status={cluster.status}
                            provider={cluster.provider}
                            confirmDeleteCallback={showConfirmDeleteModal}
                        />
                    );
                })}
            </>
        );
    }

    // Returns view to be rendered when no Unamanged clusters are present.
    function NoUnmanagedClustersSection() {
        return (
            <>
                <div
                    cds-layout="grid horizontal cols:8 p:md"
                    className="section-raised mgmt-cluster-no-cluster-container"
                    data-testid="no-clusters-messaging"
                >
                    <div cds-layout="grid horizontal cols:12 gap:lg gap@md:lg">
                        <div cds-text="title">Unmanaged Cluster not found</div>
                        <div cds-text="body">
                            Create an Unmanaged Cluster through a guided series of steps.
                            <br />
                            <br />
                            <a
                                href="https://tanzucommunityedition.io/docs/v0.12/planning/#managed-cluster"
                                target="_blank"
                                rel="noreferrer"
                                cds-text="link"
                            >
                                Learn more about Unmanaged Clusters
                            </a>
                        </div>
                    </div>
                </div>
            </>
        );
    }
}

export default UnmanagedClusterInventory;

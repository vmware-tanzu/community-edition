// React imports
import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { CdsModal, CdsModalHeader, CdsModalActions, CdsModalContent } from '@cds/react/modal';

// App imports
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import PageLoading from '../../shared/components/PageLoading/PageLoading';
import PageNotification, { Notification, NotificationStatus } from '../../shared/components/PageNotification/PageNotification';
import { UnmanagedCluster } from '../../swagger-api';
import { UnmanagedService } from '../../swagger-api/services/UnmanagedService';
import UnmanagedClusterCard from './unmanaged-cluster-common/UnmanagedClusterCard/UnmanagedClusterCard';
import './UnmanagedClusterInventory.scss';

function UnmanagedClusterInventory() {
    const [showPageLoading, setShowPageLoading] = useState<boolean>(false);
    const [notification, setNotification] = useState<Notification | null>(null);
    const [unmanagedClusters, setUnmanagedClusters] = useState<UnmanagedCluster[]>([]);
    const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
    const [clusterNameForDeletion, setClusterNameForDeletion] = useState<string>('');
    const navigate = useNavigate();

    // sets notification to null to dismiss alert
    function dismissAlert() {
        setNotification(null);
    }

    const retrieveUnmanagedClusters = async () => {
        setShowPageLoading(true);
        try {
            const data = await UnmanagedService.getUnmanagedClusters();
            data.filter((cluster) => {
                if (cluster.status === 'Unknown') {
                    setNotification({
                        status: NotificationStatus.DANGER,
                        message: `Unmanaged Clusters ${cluster.name} is failed: Please delete this cluster before attempting to create a new unmanaged cluster.`,
                    } as Notification);
                }
            });
            setUnmanagedClusters(data);
        } catch (e) {
            setNotification({
                status: NotificationStatus.DANGER,
                message: `Unable to retrieve Unmanaged Clusters: ${e}`,
            } as Notification);
        } finally {
            setShowPageLoading(false);
        }
    };

    useEffect(() => {
        retrieveUnmanagedClusters().then(() => {
            console.info('initial loading of unmanaged clusters list completed');
        });
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
        return unmanagedClusters.length > 0;
    }

    const deleteUnmanagedCluster = async (clusterName: string) => {
        setShowDeleteModal(!showDeleteModal);
        try {
            await UnmanagedService.deleteUnmanagedCluster(clusterName);
            setNotification({
                status: NotificationStatus.SUCCESS,
                message: `Successfully deleted cluster ${clusterName}`,
            } as Notification);
        } catch (e) {
            setNotification({
                status: NotificationStatus.DANGER,
                message: `Unable to delete Unmanaged Clusters ${clusterName}: ${e}`,
            } as Notification);
        }
    };

    return (
        <div className="unmanaged-cluster-landing-container" cds-layout="grid vertical col:12 gap:lg align:fill">
            <div cds-layout="grid horizontal col:12">
                <div cds-layout="vertical gap:md gap@md:lg col@sm:8 col:8">
                    {Header()}
                    <PageNotification notification={notification} closeCallback={dismissAlert}></PageNotification>
                    {MainContent()}
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
                    Unmanaged clusters are a single node, local workstation cluster suitable for a development/test environment. It requires
                    minimal local resources and is fast to deploy. It provides support for running multiple clusters. The default Tanzu
                    Community Edition package repository is automatically installed when you deploy an unmanaged cluster.{' '}
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
                            onClick={() => deleteUnmanagedCluster(clusterNameForDeletion).then(() => retrieveUnmanagedClusters())}
                            data-testid="delete-cluster-btn"
                        >
                            Delete
                        </CdsButton>
                    </CdsModalActions>
                </CdsModal>
            </>
        );
    }

    function MainContent() {
        if (showPageLoading) {
            return (
                <div cds-layout="vertical col:12 p-y:xxl">
                    <PageLoading message="Retrieving unmanaged clusters..."></PageLoading>
                </div>
            );
        } else {
            return hasUnmanagedClusters() ? UnmanagedClustersSection() : NoUnmanagedClustersSection();
        }
    }

    // Returns view to be rendered when Unmanaged clusters are present.
    function UnmanagedClustersSection() {
        return (
            <>
                <CdsButton className="cluster-action-btn" status="primary" onClick={() => navigate(NavRoutes.UNMANAGED_CLUSTER_WIZARD)}>
                    <CdsIcon shape="cluster"></CdsIcon>
                    create an unmanaged cluster
                </CdsButton>
                <div cds-text="body">The following clusters were discovered on this workstation.</div>
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

    // Returns view to be rendered when no Unmanaged clusters are present.
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
                            <p cds-layout="m-t:none">Create an Unmanaged Cluster through a guided series of steps.</p>
                            <CdsButton
                                className="cluster-action-btn"
                                status="primary"
                                onClick={() => navigate(NavRoutes.UNMANAGED_CLUSTER_WIZARD)}
                            >
                                <CdsIcon shape="cluster"></CdsIcon>
                                create an unmanaged cluster
                            </CdsButton>
                        </div>
                    </div>
                </div>
            </>
        );
    }
}

export default UnmanagedClusterInventory;

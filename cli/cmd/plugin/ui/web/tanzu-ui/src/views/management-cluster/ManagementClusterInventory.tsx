// React imports
import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { CdsModal, CdsModalActions, CdsModalContent, CdsModalHeader } from '@cds/react/modal';

// App imports
import ManagementClusterCard from './ManagementClusterCard';
import { ManagementCluster } from '../../swagger-api/models/ManagementCluster';
import { ManagementService } from '../../swagger-api';
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import './ManagementClusterInventory.scss';
import PageNotification, { Notification, NotificationStatus } from '../../shared/components/PageNotification/PageNotification';

function ManagementClusterInventory() {
    const [showLoading, setShowLoading] = useState<boolean>(false);
    const [notification, setNotification] = useState<Notification | null>(null);
    const [managementClusters, setManagementClusters] = useState<ManagementCluster[]>([]);
    const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
    const [clusterNameForDeletion, setClusterNameForDeletion] = useState<string>('');

    const navigate = useNavigate();

    // sets notification to null to dismiss alert
    function dismissAlert() {
        setNotification(null);
    }

    // Retrieve management clusters list on page load; then set loading to false
    useEffect(() => {
        retrieveManagementClusters().then(() => {
            console.info('initial loading of management clusters list completed');
        });
    }, []);

    const retrieveManagementClusters = async () => {
        setShowLoading(true);
        try {
            const data = await ManagementService.getMgmtClusters();
            setManagementClusters(data);
        } catch (e) {
            setNotification({
                status: NotificationStatus.DANGER,
                message: `Unable to retrieve Management Clusters: ${e}`,
            } as Notification);
        } finally {
            setShowLoading(false);
        }
    };

    // Helper function returns true if management clusters exist; otherwise returns false
    function hasManagementClusters() {
        return managementClusters.length ? true : false;
    }

    // Helper function to be passed to ManagementClusterCard components and leveraged for displaying confirm delete
    // modal window when user clicks Delete button. Takes cluster name as parameter for displaying this name in the
    // title and body of the modal window.
    function showConfirmDeleteModal(clusterName?: string) {
        if (!clusterName) {
            console.warn(`toggleConfirmDeleteModal expected a cluster name string; instead got ${clusterName}`);
            return;
        }
        setClusterNameForDeletion(clusterName);
        setShowDeleteModal(!showDeleteModal);
    }

    // Handler function for confirmed delete action triggered in the modal window.
    // Calls delete API with cluster name, then retrieves updated list of management clusters.
    // TODO: test e2e as management cluster deletion may not be immediate, so we may have to use another mechanism
    // for updating management cluster list periodically
    const deleteManagementCluster = async (clusterName: string) => {
        setShowDeleteModal(!showDeleteModal);
        try {
            await ManagementService.deleteMgmtCluster(clusterName);

            setNotification({
                status: NotificationStatus.SUCCESS,
                message: `Management Cluster ${clusterName} has been deleted`,
            } as Notification);
        } catch (e) {
            setNotification({
                status: NotificationStatus.DANGER,
                message: `Unable to delete Management Cluster ${clusterName}: ${e}`,
            } as Notification);
        }
    };

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
                            Deleting {clusterNameForDeletion} stops this Management Cluster and removes it from the provider you created it
                            on.
                        </p>
                    </CdsModalContent>
                    <CdsModalActions>
                        <CdsButton action="outline" onClick={() => setShowDeleteModal(false)}>
                            Cancel
                        </CdsButton>
                        <CdsButton
                            status="danger"
                            onClick={() => deleteManagementCluster(clusterNameForDeletion).then(() => retrieveManagementClusters())}
                            data-testid="delete-cluster-btn"
                        >
                            Delete
                        </CdsButton>
                    </CdsModalActions>
                </CdsModal>
            </>
        );
    }

    // Returns view to be rendered when no management clusters are present.
    function NoManagementClustersSection() {
        return (
            <>
                <div
                    cds-layout="grid horizontal cols:8 p:lg"
                    className="section-raised mgmt-cluster-no-cluster-container"
                    data-testid="no-clusters-messaging"
                >
                    <div cds-layout="grid horizontal cols:12 gap:lg gap@md:lg">
                        <div cds-text="section">Management Cluster not found</div>
                        <div cds-text="body">
                            <p cds-layout="m-t:none">
                                Create a Management Cluster on your preferred cloud provider through a guided series of steps.
                            </p>
                            <p>
                                This cluster will manage new workload clusters you create for your workloads.{' '}
                                <a
                                    href="https://tanzucommunityedition.io/docs/v0.12/planning/#managed-cluster"
                                    target="_blank"
                                    rel="noreferrer"
                                    cds-text="link"
                                >
                                    Learn more about Management Clusters
                                </a>
                                .
                            </p>
                            <CdsButton onClick={() => navigate(NavRoutes.MANAGEMENT_CLUSTER_SELECT_PROVIDER)}>
                                <CdsIcon shape="blocks-group"></CdsIcon>create a management cluster
                            </CdsButton>
                        </div>
                    </div>
                </div>
            </>
        );
    }

    // Returns view to be rendered when management clusters are present.
    function ManagementClustersSection() {
        return (
            <>
                <div cds-text="body">The following clusters were discovered on this workstation.</div>
                <div>
                    <CdsButton onClick={() => navigate(NavRoutes.MANAGEMENT_CLUSTER_SELECT_PROVIDER)}>
                        <CdsIcon shape="blocks-group"></CdsIcon>create a management cluster
                    </CdsButton>
                </div>
                {managementClusters.map((cluster: ManagementCluster) => {
                    return (
                        <ManagementClusterCard
                            key={cluster.name}
                            name={cluster.name}
                            path={cluster.path}
                            context={cluster.context}
                            confirmDeleteCallback={showConfirmDeleteModal}
                        ></ManagementClusterCard>
                    );
                })}
            </>
        );
    }

    return (
        <>
            <div className="management-cluster-landing-container" cds-layout="vertical gap:md col@sm:12 grid">
                <div cds-layout="vertical col:8 gap:lg">
                    <div cds-text="title" cds-layout="horizontal align:vertical-center">
                        <CdsIcon cds-layout="m-r:sm" shape="blocks-group" size="lg"></CdsIcon>
                        Management Clusters
                    </div>
                    <PageNotification notification={notification} closeCallback={dismissAlert}></PageNotification>
                    {hasManagementClusters() ? ManagementClustersSection() : NoManagementClustersSection()}
                </div>
                <div cds-layout="col:4" className="mgmt-cluster-admins-img"></div>
                {renderConfirmDeleteModal()}
            </div>
        </>
    );
}

export default ManagementClusterInventory;
